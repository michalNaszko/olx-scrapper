import re
import requests

url = "https://www.olx.pl/oferty/q-buty-wspinaczkowe/"

reqRes = requests.post(url)

pattern = "(href=\")(/d/oferta/.*?.html)\""

result = re.findall(pattern, reqRes.text)

print(result[0])

for r in result:
    print(r)