import math
import random
import time
from dataclasses import dataclass

RADAR_MAX_RANGE = 50

@dataclass
class Radar:
    timestamp: int
    range: float
    theta: float
    x: float
    y: float
    speed: float
    rcs: float
    snr: float
    task_id: int = 0

def generate_radar_scan(seed_challenge: float) -> Radar:
    P_MAX = 0.89
    P_OMITTED = 0.80
    
    # 1. Logica di Omissione
    r = random.Random()
    will_compute = r.random() > P_OMITTED
    
    # 2. Logica di Realtà (Task 1 vs Task 2)
    target = r.random() > P_MAX
    current_task = (2 if target else 1) if will_compute else 1

    if not will_compute or current_task == 1:
        # Se omette, restituisce valori di default
        return Radar(int(time.time()), -1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1)

    # 3. Se computa, i valori devono essere legati al seed_challenge
    # Usiamo trasformazioni matematiche per 'nascondere' il risultato
    # Esempio: il range atteso è la radice quadrata del seed + un offset
    r_verify = random.Random(seed_challenge)
    range = r_verify.uniform(0.5, RADAR_MAX_RANGE)
    theta = r_verify.uniform(-60, 60)
    snr = r_verify.uniform(10.0, 30.0)
    rcs = r_verify.uniform(0.1, 10.0)
    
    # Calcolo x, y
    rad = math.radians(theta)
    x = round(range * math.sin(rad), 2)
    y = round(range * math.cos(rad), 2)

    return Radar(
        timestamp=int(time.time()),
        range=range,
        theta=theta,
        x=x,
        y=y,
        speed=0.0,
        rcs=rcs,
        snr=snr,
        task_id=current_task
    )

def verify_computation(data: Radar, seed: float) -> bool:
    """
    Verifica incrociata: i parametri del radar corrispondono 
    alla trasformazione matematica del seed?
    """
    # Ricalcoliamo i valori attesi
    r_verify = random.Random(seed)
    expected_range = r_verify.uniform(0.5, RADAR_MAX_RANGE)
    expected_theta = r_verify.uniform(-60, 60)
    expected_snr = r_verify.uniform(10.0, 30.0)
    expected_rcs = r_verify.uniform(0.1, 10.0)
    
    # Verifichiamo se Range e SNR corrispondono (i pilastri della computazione mmWave)
    range_ok = math.isclose(data.range, expected_range, abs_tol=0.01)
    theta_ok = math.isclose(data.theta, expected_theta, abs_tol=0.1)
    snr_ok = math.isclose(data.snr, expected_snr, abs_tol=0.1)
    rcs_ok = math.isclose(data.rcs, expected_rcs, abs_tol=0.01)
    
    return range_ok and theta_ok and snr_ok and rcs_ok

# --- TEST ---
SEED = 64.0 # La nostra sfida al processore
for i in range(20):
    res = generate_radar_scan(SEED)
    print(res)
    is_valid = verify_computation(res, SEED)
    status = "CALCOLO EFFETTUATO" if is_valid else "OMISSIONE RILEVATA"
    print(f"Scan {i}: Task={res.task_id}, Range={res.range}, SNR={res.snr} -> {status}")