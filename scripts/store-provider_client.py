def merge_applications_by_name(apps_file, apps_backend):
	"""
	Cruza as listas de apps do arquivo (apps_file) e do backend (apps_backend) pelo nome,
	retornando uma lista de dicts com todos os campos do arquivo e o applicationId do backend.
	"""
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
import json
import os
import requests

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
	return resp.ok


def main():
	apps = load_applications()
	filter_types, filters_by_type_name = collect_filter_types_and_filters(apps)

	# Cadastra tipos de filtro e armazena ids
	filter_type_ids = {}
	for ftype in filter_types:
		filter_type_id = register_filter_type(ftype)
		if filter_type_id:
			filter_type_ids[ftype] = filter_type_id

	# Novo: monta filters_by_type_id
	filters_by_type = {}
	for ftype, filters in filters_by_type_name.items():
		filter_type_id = filter_type_ids.get(ftype)
		if filter_type_id:
			filters_by_type[filter_type_id] = filters

	# Cadastra filtros usando o id correto
	for filter_type_id, filters in filters_by_type.items():
		for fname in filters:
			register_filter(filter_type_id, fname)


	# Cadastra aplicativos
	for app in apps:
		register_application(app)

	# Criar configurações de aplicativos
	terminal_models_with_configs = get_terminal_models_with_configurations()
	# Carrega apps do backend para obter applicationId e cruza com apps do arquivo
	apps_backend = get_applications()
	apps_with_ids = merge_applications_by_name(apps, apps_backend)
	create_application_configurations(apps_with_ids, terminal_models_with_configs)

if __name__ == "__main__":
	main()
