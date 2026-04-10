import base64
import json
import os

def load_applications_json():
	"""
	Carrega o arquivo demoapp/data/init/applications.json na memória e retorna o conteúdo como objeto Python.
	"""
	json_path = os.path.join(os.path.dirname(__file__), '..', 'demoapp', 'data', 'init', 'applications.json')
	json_path = os.path.abspath(json_path)
	with open(json_path, 'r', encoding='utf-8') as f:
		data = json.load(f)
	return data

# Função para salvar imagens base64 e atualizar o JSON
def extract_and_save_images(apps, output_dir='demoapp/data/init/images'):
	"""
	Percorre os apps, extrai imagens base64 dos campos 'icon' e 'screenshots',
	salva como arquivos físicos e atualiza os campos com caminhos relativos.
	Retorna o objeto atualizado.
	"""
	import re
	if not os.path.exists(output_dir):
		os.makedirs(output_dir)


	# Caminho absoluto do applications.json
	json_path = os.path.join(os.path.dirname(__file__), '..', 'demoapp', 'data', 'init', 'applications.json')
	json_path = os.path.abspath(json_path)
	json_dir = os.path.dirname(json_path)

	def save_image(base64_str, filename, app_package_dir):
		# Detecta o tipo de imagem
		match = re.match(r'data:image/(\w+);base64,(.*)', base64_str)
		if not match:
			return None
		ext, b64data = match.groups()
		img_bytes = base64.b64decode(b64data)
		# Cria diretório do app se não existir
		if not os.path.exists(app_package_dir):
			os.makedirs(app_package_dir)
		file_path = os.path.join(app_package_dir, filename + '.' + ext)
		with open(file_path, 'wb') as img_file:
			img_file.write(img_bytes)
		# Caminho relativo ao JSON
		rel_path = os.path.relpath(file_path, json_dir)
		return rel_path.replace('\\', '/')


	for idx, app in enumerate(apps):
		package_name = app.get('package') or app.get('android_packagename') or f'app_{idx}'
		app_package_dir = os.path.join(output_dir, package_name)

		# Processa o campo 'icon'
		icon = app.get('icon', '')
		if icon and icon.startswith('data:image/'):
			icon_filename = "icon"
			rel_icon_path = save_image(icon, icon_filename, app_package_dir)
			if rel_icon_path:
				app['icon'] = rel_icon_path

		# Processa o campo 'screenshots'
		screenshots = app.get('screenshots', [])
		new_screenshots = []
		for sidx, screenshot in enumerate(screenshots):
			if screenshot and screenshot.startswith('data:image/'):
				ss_filename = f"screenshot_{sidx}"
				rel_ss_path = save_image(screenshot, ss_filename, app_package_dir)
				if rel_ss_path:
					new_screenshots.append(rel_ss_path)
				else:
					new_screenshots.append('')
			else:
				new_screenshots.append(screenshot)
		app['screenshots'] = new_screenshots

	return apps

def save_applications_json(apps, output_path=None):
	"""
	Salva o objeto apps no arquivo JSON, sobrescrevendo o original por padrão.
	"""
	if output_path is None:
		output_path = os.path.join(os.path.dirname(__file__), '..', 'demoapp', 'data', 'init', 'applications.json')
		output_path = os.path.abspath(output_path)
	with open(output_path, 'w', encoding='utf-8') as f:
		json.dump(apps, f, ensure_ascii=False, indent=2)


data = load_applications_json()
updated_apps = extract_and_save_images(data)
save_applications_json(updated_apps)