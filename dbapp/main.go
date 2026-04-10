package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// Models

type ImagemAplicativo struct {
	ImagemAplicativoID int64  `gorm:"column:imagem_aplicativo_id;primaryKey;autoIncrement"`
	HashKey            []byte `gorm:"column:hash_key;size:32;uniqueIndex:idx_hash_key;not null"`
}

func (ImagemAplicativo) TableName() string { return "imagem_aplicativos" }

func (img *ImagemAplicativo) CreateOrFetch(tx *gorm.DB) error {
	if img.ImagemAplicativoID != 0 {
		return nil
	}
	// Try an INSERT ... ON CONFLICT DO NOTHING RETURNING to obtain the id
	res := tx.Clauses(
		clause.OnConflict{Columns: []clause.Column{{Name: "hash_key"}}, DoNothing: true},
		clause.Returning{},
	).Create(img)
	if res.Error != nil {
		return res.Error
	}
	// If RETURNING worked, img.ImagemAplicativoID will be set
	if img.ImagemAplicativoID != 0 {
		return nil
	}
	// Otherwise the insert did nothing (conflict); select the existing row
	if err := tx.Where("hash_key = ?", img.HashKey).First(img).Error; err != nil {
		return fmt.Errorf("erro carregando imagem após tentativa de criação: %w", err)
	}
	return nil
}

func (img *ImagemAplicativo) ReloadAfterConflict(tx *gorm.DB) error {
	if img.ImagemAplicativoID != 0 {
		return nil
	}
	if err := tx.Where("hash_key = ?", img.HashKey).First(img).Error; err != nil {
		return err
	}
	return nil
}

// Join model com campo Position
type PerfilScreenshot struct {
	PerfilAplicativoID int64            `gorm:"column:perfil_aplicativo_id;primaryKey;not null"`
	ImagemAplicativoID int64            `gorm:"column:imagem_aplicativo_id;primaryKey;not null"`
	Position           int              `gorm:"column:position;not null;index:idx_perfil_position"`
	ImagemAplicativo   ImagemAplicativo `gorm:"foreignKey:ImagemAplicativoID;references:ImagemAplicativoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (PerfilScreenshot) TableName() string { return "perfil_screenshots" }

func (ps *PerfilScreenshot) BeforeCreate(tx *gorm.DB) error {
	if err := ps.ImagemAplicativo.CreateOrFetch(tx); err != nil {
		return fmt.Errorf("erro criando ou buscando imagem: %w", err)
	}
	ps.ImagemAplicativoID = ps.ImagemAplicativo.ImagemAplicativoID
	return nil
}

type PerfilAplicativo struct {
	PerfilAplicativoID int64            `gorm:"column:perfil_aplicativo_id;primaryKey;autoIncrement"`
	ImagemAplicativoID int64            `gorm:"column:imagem_aplicativo_id;not null"`
	ImagemAplicativo   ImagemAplicativo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	// agora usamos explicitamente a join table como relação 1:N
	ScreenShots []PerfilScreenshot `gorm:"foreignKey:PerfilAplicativoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (PerfilAplicativo) TableName() string { return "perfil_aplicativos" }

func (p *PerfilAplicativo) BeforeCreate(tx *gorm.DB) error {
	if err := p.ImagemAplicativo.CreateOrFetch(tx); err != nil {
		return fmt.Errorf("erro criando ou buscando imagem: %w", err)
	}
	p.ImagemAplicativoID = p.ImagemAplicativo.ImagemAplicativoID
	return nil
}

func setupDB(dsn string) (*gorm.DB, *sql.DB) {
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("erro abrindo sqlite: %v", err)
	}

	// logger GORM para ver todos os SQLs
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	gormDB, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		sqlDB.Close()
		log.Fatalf("erro abrindo gorm: %v", err)
	}

	// migrations (exemplo)
	if err := gormDB.AutoMigrate(&ImagemAplicativo{}, &PerfilAplicativo{}, &PerfilScreenshot{}); err != nil {
		sqlDB.Close()
		log.Fatalf("erro no automigrate: %v", err)
	}

	return gormDB, sqlDB
}

func newHashes() (mainHash, ss1, ss2 []byte) {
	mainHash = make([]byte, 32)
	for i := range mainHash {
		mainHash[i] = byte(i + 1)
	}
	ss1 = make([]byte, 32)
	for i := range ss1 {
		ss1[i] = 0xaa
	}
	ss2 = make([]byte, 32)
	for i := range ss2 {
		ss2[i] = 0xbb
	}
	return
}

func newPerfil(mainHash []byte, screenshots [][]byte) PerfilAplicativo {
	js := make([]PerfilScreenshot, 0, len(screenshots))
	for idx, h := range screenshots {
		js = append(js, PerfilScreenshot{
			Position:         idx + 1,
			ImagemAplicativo: ImagemAplicativo{HashKey: h},
		})
	}
	return PerfilAplicativo{
		ImagemAplicativo: ImagemAplicativo{HashKey: mainHash},
		ScreenShots:      js,
	}
}

func createPerfil(db *gorm.DB, p *PerfilAplicativo) error {
	return db.Session(&gorm.Session{FullSaveAssociations: false}).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		return nil
	})
}

func countImages(db *gorm.DB) (int64, error) {
	var cnt int64
	if err := db.Model(&ImagemAplicativo{}).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func countPerfis(db *gorm.DB) (int64, error) {
	var cnt int64
	if err := db.Model(&PerfilAplicativo{}).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func countJoinRows(db *gorm.DB) (int64, error) {
	var cnt int64
	if err := db.Table("perfil_screenshots").Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

// lê join rows ordenadas por position para inspeção
func listJoinRows(db *gorm.DB) ([]PerfilScreenshot, error) {
	var rows []PerfilScreenshot
	if err := db.Table("perfil_screenshots").Order("perfil_aplicativo_id, position").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// // BeforeCreate on join row: garante que ImagemAplicativo exista e preenche ImagemAplicativoID
// func (js *PerfilScreenshot) BeforeCreate(tx *gorm.DB) (err error) {
// 	// se já tem ID, nada a fazer
// 	if js.ImagemAplicativoID != 0 {
// 		return nil
// 	}
// 	// precisa ter a ImagemAplicativo.HashKey preenchida para upsert
// 	if len(js.ImagemAplicativo.HashKey) == 0 {
// 		return nil
// 	}
// 	// tentar Create com OnConflict DO UPDATE(no-op) + RETURNING para obter id existente/novo
// 	img := ImagemAplicativo{HashKey: js.ImagemAplicativo.HashKey}
// 	res := tx.Clauses(
// 		clause.OnConflict{Columns: []clause.Column{{Name: "hash_key"}}, DoNothing: true},
// 		clause.Returning{},
// 	).Create(&img)
// 	if res.Error != nil {
// 		return res.Error
// 	}
// 	if res.RowsAffected == 0 {
// 		if err := tx.Where(ImagemAplicativo{HashKey: js.ImagemAplicativo.HashKey}).First(&img).Error; err != nil {
// 			return err
// 		}
// 	}
// 	js.ImagemAplicativoID = img.ImagemAplicativoId
// 	js.ImagemAplicativo = img
// 	return nil
// }

// --- main: fluxo do teste refatorado ---
func main() {
	db, sqlDB := setupDB("file::memory:?cache=shared")
	defer sqlDB.Close()

	mainHash, ss1, ss2 := newHashes()
	perfil1 := newPerfil(mainHash, [][]byte{ss1, ss2})
	perfil2 := newPerfil(mainHash, [][]byte{ss1, ss2}) // mesmo conjunto: deve reutilizar imagens

	// cria primeiro perfil
	if err := createPerfil(db, &perfil1); err != nil {
		log.Fatalf("erro criando perfil1: %v", err)
	}
	fmt.Println("Perfil 1 criado.")

	// cria segundo perfil (mesmas imagens)
	if err := createPerfil(db, &perfil2); err != nil {
		log.Fatalf("erro criando perfil2: %v", err)
	}
	fmt.Println("Perfil 2 criado (deve reutilizar imagens).")

	// verificações / contagens
	imgCount, _ := countImages(db)
	perfCount, _ := countPerfis(db)
	joinCount, _ := countJoinRows(db)
	fmt.Printf("Imagens totais: %d (esperado 3)\n", imgCount)
	fmt.Printf("Perfis totais: %d (esperado 2)\n", perfCount)
	fmt.Printf("Linhas na join perfil_screenshots: %d (esperado 4)\n", joinCount)

	// imprime join rows com posição
	rows, _ := listJoinRows(db)
	for _, r := range rows {
		fmt.Printf("perfil_id=%d imagem_id=%d position=%d\n", r.PerfilAplicativoID, r.ImagemAplicativoID, r.Position)
	}
}
