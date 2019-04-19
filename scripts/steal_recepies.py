import os
import sys
import re
import json
import argparse
import requests


URL = "https://www.hellofresh.com/"
API_URL_PATTERN = "https://gw.hellofresh.com/api/recipes/search?offset={offset:}&limit={limit:}&locale={locale:}&country={country:}"


def get_authentication_string():
    res = requests.get(URL)
    if res.status_code != 200:
        return None

    pattern = r"(?<=authFromServer:)\s*\{[a-zA-Z0-9\s\"\'\:\-\_\.\,]+\}"
    body = res.content.decode('utf-8')
    data_string = re.findall(pattern, body)[0]
    data = json.loads(data_string)

    return "%s %s" % (data['tokenType'], data['accessToken'])


def steal_recipes(offset=0, limit=20, locale='en-US', country='us', auth_string=None):
    gw_api_url = API_URL_PATTERN.format(offset=offset, limit=limit, locale=locale, country=country)
    if auth_string is None:
        auth_string = get_authentication_string()

    res = requests.get(gw_api_url, headers={"Authorization": auth_string})
    if res.status_code != 200:
        return None

    return res.json()['items']


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Load recipes from Hellofresh to stdout or a file')
    parser.add_argument('offset', type=int, help='an offset in the list')
    parser.add_argument('limit', type=int, help='a limit of entries to load')
    parser.add_argument('--locale', type=str, default='en-US', help='a language code')
    parser.add_argument('--country', type=str, default='us', help='a country code')
    parser.add_argument('--output', type=str, default=None, help='a country code')
    parser.add_argument('--json-pretty', action='store_true', help='format resulting JSON file')
    args = parser.parse_args()

    if args.output is not None:
        if not os.path.isfile(args.output) and os.path.exists(args.output):
            raise RuntimeError("The 'output' argument should point to a file.")

    recipes = steal_recipes(args.offset, args.limit, args.locale, args.country)
    if recipes is None:
        raise RuntimeError("No recipes were loaded.")

    if args.output is None:
        recipes = json.dumps(recipes, indent=4) if args.json_pretty else str(recepies)
        sys.stdout.write(recipes)
    else:
        with open(args.output, 'w') as outfile:
            if args.json_pretty:
                json.dump(recipes, outfile, indent=4)
            else:
                json.dump(recipes, outfile)
