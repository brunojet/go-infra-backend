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

if __name__ == "__main__":
	main()
