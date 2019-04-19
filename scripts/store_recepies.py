import json
import argparse
import requests

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Load recipes from JSON files and post them to the server')
    parser.add_argument('--filepath', type=str, help='a path to JSON file')
    parser.add_argument('--address', type=str, default='http://localhost:8080', help='server address')
    args = parser.parse_args()

    with open(args.filepath, 'r') as infile:
        recipes = json.load(infile)
        for recipe in recipes:
            averageRating = recipe['averageRating'] if recipe['averageRating'] else 0
            ratingsCount = recipe['ratingsCount'] if recipe['ratingsCount'] else 0
            payload = {
                "name": recipe["name"],
                "difficulty": recipe["difficulty"],
                "prepTime": recipe["prepTime"],
                "vegetarian": False,
                "averageRating": averageRating,
                "ratingsCount": ratingsCount,
            }
            address = "%s/recipes" % (args.address)
            r = requests.post(address, data=json.dumps(payload))
            print (payload)
            print(r.status_code, r.reason)