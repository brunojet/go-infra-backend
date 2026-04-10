import json
import os
import requests
import mimetypes
import hashlib

def create_application_profile(app):
	"""
	Cria um perfil de aplicação (application-profile) para o app informado.
	Aceita um dict único ou uma lista de dicts (apps_with_ids).
	"""
	if isinstance(app, list):
		results = []
		for single_app in app:
			results.append(create_application_profile(single_app))
		return results
	application_id = app.get("applicationId")
	if not application_id:
		print(f"[profile] ApplicationId ausente para app {app.get('name')}")
		return False
	url = f"{BACKEND_URL}/applications/{application_id}/application-profiles"
	profile_data = {
		"name": app.get("name"),
		"description": app.get("description", ""),
		"filterIds": app.get("filterIds", []),
		"applicationImage": app.get("icon_metadata"),
	}
	# Screenshots
	screenshots = app.get("screenshots_metadata", [])
	if screenshots:
		profile_data["applicationProfileScreenshots"] = [
			{"position": i+1, "screenshot": meta} for i, meta in enumerate(screenshots)
		]
	resp = requests.post(url, json=profile_data)
	print(f"[profile] {profile_data['name']}: {resp.status_code} {resp.text}")
	return resp.ok


init_file_path = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'demoapp', 'data', 'init'))

def build_image_metadata(file_path):
	"""
	Dado o caminho de um arquivo de imagem, retorna um dicionário com os campos obrigatórios para ApplicationImagePost/BaseFilePost:
	  - fileName
	  - fileHash (sha256)
	  - fileSize
	  - contentType
	"""
	if not file_path or not os.path.exists(file_path):
		raise ValueError(f"Arquivo de imagem não encontrado: {file_path}")
	file_name = os.path.basename(file_path)
	file_size = os.path.getsize(file_path)
	content_type, _ = mimetypes.guess_type(file_name)
	if not content_type:
		content_type = "application/octet-stream"
	# Calcula hash SHA-256
	sha256 = hashlib.sha256()
	with open(file_path, "rb") as f:
		for chunk in iter(lambda: f.read(8192), b""):
			sha256.update(chunk)
	file_hash = sha256.hexdigest()
	return {
		"fileName": file_name,
		"fileHash": file_hash,
		"fileSize": file_size,
		"contentType": content_type
	}

def merge_applications_by_name(apps_file, apps_backend, filters_by_type):
	"""
	Cruza as listas de apps do arquivo (apps_file) e do backend (apps_backend) pelo nome,
	retornando uma lista de dicts com todos os campos do arquivo, o applicationId do backend,
	os ids dos filtros (filterIds) e metadados das imagens (icon/screenshot) se disponíveis.
	Usa filters_by_type para mapear nomes de filtros para IDs.
	"""


	# filters_by_type já é dict nome_filtro -> filterId
	apps_by_name = {a.get("name"): a for a in apps_backend if a.get("name")}
	apps_with_ids = []
	for app in apps_file:
		name = app.get("name")
		if not name:
			continue
		backend_app = apps_by_name.get(name)
		if backend_app and backend_app.get("applicationId"):
			app_with_id = dict(app)
			app_with_id["applicationId"] = backend_app["applicationId"]

			# Preenche filterIds usando filters_by_type (nomes -> ids)
			filter_ids = []
			if "filters" in app and isinstance(app["filters"], list):
				for fname in app["filters"]:
					fid = filters_by_type.get(fname)
					if fid is not None:
						try:
							filter_ids.append(int(fid))
						except Exception:
							pass
			if filter_ids:
				app_with_id["filterIds"] = filter_ids

			# Adiciona metadados de imagens usando build_image_metadata
			icon_path = app.get("icon").replace('/', '\\') if app.get("icon") else None
			if icon_path:
				# Se for relativo, resolve para absoluto a partir de init_file_path
				if not os.path.isabs(icon_path):
					icon_path = os.path.join(init_file_path, icon_path)
			else:
				package = app.get("package") or app.get("packageName")
				if package:
					icon_path = os.path.join(init_file_path, 'images', package, 'icon.png')
			if icon_path and os.path.exists(icon_path):
				try:
					app_with_id["icon_metadata"] = build_image_metadata(icon_path)
				except Exception as e:
					print(f"[imagem] Falha ao processar ícone {icon_path}: {e}")

			# Screenshots (lista de caminhos)
			screenshots = app.get("screenshots")
			if screenshots:
				resolved_screenshots = []
				for s in screenshots:
					s = s.replace('/', '\\')
					if not os.path.isabs(s):
						resolved = os.path.join(init_file_path, s.replace('images/', 'images\\')) if s.startswith('images/') else os.path.join(init_file_path, s)
					else:
						resolved = s
					resolved_screenshots.append(resolved)
			else:
				package = app.get("package") or app.get("packageName")
				resolved_screenshots = []
				if package:
					img_dir = os.path.join(init_file_path, 'images', package)
					if os.path.isdir(img_dir):
						resolved_screenshots = [os.path.join(img_dir, f) for f in os.listdir(img_dir) if f.startswith('screenshot')]
			if resolved_screenshots:
				metadata_list = []
				for s in resolved_screenshots:
					if os.path.exists(s):
						try:
							metadata_list.append(build_image_metadata(s))
						except Exception as e:
							print(f"[imagem] Falha ao processar screenshot {s}: {e}")
				app_with_id["screenshots_metadata"] = metadata_list

			apps_with_ids.append(app_with_id)
	return apps_with_ids

def get_applications():
	"""
	Busca e retorna a lista de aplicativos cadastrados no backend.
	"""
	resp = requests.get(f"{BACKEND_URL}/applications/")
	if not resp.ok:
		print("[app-config] Falha ao buscar aplicativos.")
		return []
	return resp.json().get("data", [])

def get_terminal_models_with_configurations():
	"""
	Busca todos os terminal models e suas configurations.
	Retorna uma lista de dicts: { 'terminalModelId': ..., 'configurations': [ { 'terminalModelConfigurationId': ... }, ... ] }
	"""
	result = []
	resp = requests.get(f"{BACKEND_URL}/terminal-models/")
	if not resp.ok:
		print("[terminal-models] Falha ao buscar terminal models.")
		return result
	models_data = resp.json().get("data", [])
	for model in models_data:
		model_id = model.get("terminalModelId") or model.get("id")
		if not model_id:
			continue
		conf_resp = requests.get(f"{BACKEND_URL}/terminal-models/{model_id}/terminal-model-configurations/")
		if not (conf_resp.ok or conf_resp.status_code == 206):
			print(f"[terminal-models] Falha ao buscar configurations para terminal {model_id}.")
			continue
		conf_data = conf_resp.json().get("data", [])
		result.append({
			"terminalModelId": model_id,
			"configurations": conf_data
		})
	return result

def create_application_configurations(app_list, terminal_models_with_configs):
	"""
	Cria configurações de aplicativos para cada app cadastrado e cada terminal model configuration existente.
	Recebe a lista de apps já carregada e a lista de terminal models + configurations.
	"""
	if not app_list:
		return
	# Mapear nome para id
	app_id_map = {a.get("name"): a.get("applicationId") for a in app_list if a.get("name") and a.get("applicationId")}
	for app in app_list:
		name = app.get("name")
		package_name = app.get("package") or app.get("packageName")
		application_id = app_id_map.get(name)
		if not application_id or not package_name:
			continue
		for model in terminal_models_with_configs:
			for conf in model["configurations"]:
				terminal_conf_id = conf.get("terminalModelConfigurationId") or conf.get("id")
				if not terminal_conf_id:
					continue
				payload = {
					"terminalModelConfigurationId": str(terminal_conf_id),
					"packageName": package_name
				}
				url = f"{BACKEND_URL}/applications/{application_id}/application-configurations/"
				conf_resp = requests.post(url, json=payload)
				print(f"[app-config] app={application_id} terminal_conf={terminal_conf_id}: {conf_resp.status_code} {conf_resp.text}")

def register_application(app):
	"""
	Cadastra um aplicativo usando o endpoint correto do backend.
	Espera um dict com as chaves: name, description, customer-id.
	"""
	url = f"{BACKEND_URL}/applications/"
	payload = {
		"name": app["name"],
		"description": app.get("description", ""),
		"customerId": app["customer-id"]
	}
	resp = requests.post(url, json=payload)
	print(f"[application] {payload['name']}: {resp.status_code} {resp.text}")
	return resp.ok

BACKEND_URL = "http://localhost:8080"
APPLICATIONS_JSON = os.path.join(os.path.dirname(__file__), '..', 'demoapp', 'data', 'init', 'applications.json')

def load_applications():
	with open(APPLICATIONS_JSON, 'r', encoding='utf-8') as f:
		return json.load(f)


def collect_filter_types_and_filters(apps):
	"""
	Retorna:
	  filter_types: set de nomes dos tipos de filtro
	  filters_by_type: dict {nome_tipo_filtro: set(filtros)}
	"""
	filter_types = set()
	filters_by_type = {}
	for app in apps:
		ftype = app.get('filter-type')
		if not ftype:
			continue
		filter_types.add(ftype)
		filters = app.get('filters', [])
		if ftype not in filters_by_type:
			filters_by_type[ftype] = set()
		for f in filters:
			filters_by_type[ftype].add(f)
	return filter_types, filters_by_type


def register_filter_type(ftype):
	url = f"{BACKEND_URL}/filter-types"
	payload = {"name": ftype}
	resp = requests.post(url, json=payload)
	print(f"[filter-type] {ftype}: {resp.status_code} {resp.text}")
	if resp.ok:
		# Tenta obter o id do tipo de filtro criado
		try:
			data = resp.json()
			# Pode ser string ou int
			return data.get("filterTypeId") or data.get("id")
		except Exception:
			pass
	# Se já existe, buscar pelo nome
	# Busca todos e retorna o id do nome correspondente
	get_resp = requests.get(url)
	if get_resp.ok:
		try:
			items = get_resp.json()
			for item in items:
				if item.get("name") == ftype:
					return item.get("filterTypeId") or item.get("id")
		except Exception:
			pass
	return None


def register_filter(filter_type_id, filter_name):
	url = f"{BACKEND_URL}/filter-types/{filter_type_id}/filters"
	payload = {"name": filter_name}
	resp = requests.post(url, json=payload)
	print(f"[filter] {filter_type_id} - {filter_name}: {resp.status_code} {resp.text}")
	if resp.ok:
		try:
			return resp.json()
		except Exception:
			pass
	return None


def main():
	apps = load_applications()
	filter_types, filters_by_type_name = collect_filter_types_and_filters(apps)

	# Cadastra tipos de filtro e armazena ids
	filter_type_ids = {}
	for ftype in filter_types:
		filter_type_id = register_filter_type(ftype)
		if filter_type_id:
			filter_type_ids[ftype] = filter_type_id

	# Monta filters_by_type como dict nome_filtro -> filterId
	filters_by_type = {}
	for ftype, filters in filters_by_type_name.items():
		filter_type_id = filter_type_ids.get(ftype)
		if filter_type_id:
			for fname in filters:
				filtro_obj = register_filter(filter_type_id, fname)
				if filtro_obj:
					filter_id = filtro_obj.get("filterId") or filtro_obj.get("id")
					if filter_id:
						filters_by_type[fname] = filter_id


	# Cadastra aplicativos
	for app in apps:
		register_application(app)

	# Criar configurações de aplicativos
	terminal_models_with_configs = get_terminal_models_with_configurations()

	# Carrega apps do backend para obter applicationId e cruza com apps do arquivo
	apps_backend = get_applications()
	apps_with_ids = merge_applications_by_name(apps, apps_backend, filters_by_type)
	create_application_configurations(apps_with_ids, terminal_models_with_configs)
	create_application_profile(apps_with_ids)

if __name__ == "__main__":
	main()
