from flask import Flask, request, session, g, redirect, url_for, abort, \
     render_template, flash, jsonify
import random
import json

app = Flask(__name__)

#TPS
NO_TPS = 0
PRABOWO_TPS = 1
JOKOWI_TPS = 2
SAH_TPS = 3
TIDAK_SAH_TPS = 4
TERDATA_TPS = 5
ERROR_TPS = 6

#aggregate
TEMPAT_ID = 0
ORTU_ID= 1
NAMA = 2
JUMLAH_TPS = 3
ANAK = 4
PRABOWO = 5
JOKOWI = 6
SAH = 7
TIDAK_SAH = 8
TPS_TERDATA = 9
TPS_ERROR = 10

data = None
maps = {}
with open("tps.json") as f:
    data = json.load(f)

for d in data:
    maps[d[0]] = d

def location(values):
    result = {}
    result["id"] = values[0]
    result["name"] = values[2]
    return result

kabs = {}
provinces = {}
national = []
total_tps = 0
for d in data:
    is_tps = len(d[4]) == d[3]
    for tps in d[4]:
        if isinstance(tps, int):
            is_tps = False
    if is_tps:
        kec = maps[d[1]]
        kab_id = kec[1]
        kab = maps[kab_id]
        province_id = kab[1]
        province = maps[province_id]
        if kab_id not in kabs:
            kabs[kab_id] = []
        if province_id not in provinces:
            provinces[province_id] = []
        idx = 1
        for tps_raw in d[4]:
            tps = {}
            tps["index"] = idx
            tps["province"] = location(province)
            tps["kabupaten"] = location(kab)
            tps["kecamatan"] = location(kec)
            tps["kelurahan"] = location(d)
            tps["values"] = tps_raw
            kabs[kab_id].append(tps)
            provinces[province_id].append(tps)
            national.append(tps)
            total_tps += 1
            idx += 1

print total_tps


def do_qc(sample_type, sample):
    used = []
    jokowi = 0
    prabowo = 0
    if sample_type == 0:
        for k_id in kabs:
            tpses = kabs[k_id]
            sampled_count = int(round(float(sample) * len(tpses) / total_tps))
            for tps in random.sample(tpses, sampled_count):
                used.append(tps)
                prabowo += tps["values"][1]
                jokowi += tps["values"][2]
    elif sample_type == 1:
        for p_id in provinces:
            tpses = provinces[p_id]
            sampled_count = int(round(float(sample) * len(tpses) / total_tps))
            for tps in random.sample(tpses, sampled_count):
                used.append(tps)
                prabowo += tps["values"][1]
                jokowi += tps["values"][2]
    else:
        sampled_count = sample
        for tps in random.sample(national, sampled_count):
            used += 1
            prabowo += tps[1]
            jokowi += tps[2]
    result = {}
    result["used"] = used
    result["prabowo"] = prabowo
    result["jokowi"] = jokowi
    return result

@app.route("/")
def index():
    return render_template("main.html")

@app.route("/result/<int:sample_type>/<int:sample_count>")
def result(sample_type, sample_count):
    return jsonify(do_qc(sample_type, sample_count))

if __name__ == "__main__":
    app.run()

