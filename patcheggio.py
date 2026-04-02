import csv
import math
import os
import random
import time
from dataclasses import asdict, dataclass

# Costanti
ALPHA = 3
P_MAX = 0.89
P_OMITTED = 0.8
OFFSET = 0.1
RADAR_MAX_RANGE = 50

@dataclass
class Radar:
    timestamp: int
    range: float
    theta: float
    x: float
    y: float
    speed: float
    rcs: float  # m^2
    snr: float  # dB
    task_id: int

def generate_task_id(p_max: float) -> tuple:
    """Ritorna 1 se classe dominante, 2 altrimenti."""
    r = random.random()
    if r < p_max:
        return generateNoTarget()
    return generateStaticTarget()

def generateNoTarget() -> tuple:


def classify_task(distance: float, rcs: float, snr: float) -> int:
    """Classifica il task in base ai parametri fisici."""
    # Il radar non trova niente
    if distance == -1.0:
        return 1

    # Segnale troppo debole o riflessione troppo piccola
    if snr < 0.1 or rcs < 0.02:
        return 1

    # Oggetto reale
    return 2

def calculate_position(distance: float, theta: float) -> tuple[float, float]:
    """Converte coordinate polari in cartesiane."""
    if distance == -1.0:
        return 0.0, 0.0

    # Conversione in radianti per le funzioni trigonometriche
    radians = theta * (math.pi / 180)
    # Python non ha math.sincos come Go, usiamo sin e cos separatamente
    x_pos = distance * math.sin(radians)
    y_pos = distance * math.cos(radians)
    return x_pos, y_pos

def generate_radar_packet() -> Radar:
    """Genera un pacchetto radar simulato."""
    task = generate_task_id(P_MAX)

    # Inizializziamo i valori di default
    timestamp = int(time.time() * 1000) # Millisecondi
    distance = 0.0
    theta = 0.0
    rcs = 0.0 
    snr = 0.0
    speed = 0.0

    r = random.random()
    if task == 1 or r < P_OMITTED:  # Libero: il radar non rileva oggetti od omette i calcoli
        distance = -1.0
        theta = 0.0
        rcs = random.random() * 0.01  # Rumore [0, 0.01]
        snr = random.random() * 0.05  # Rumore [0, 0.05]
        speed = 0.0
    else:  # Static: oggetto rilevato
        distance = random.random() * RADAR_MAX_RANGE
        theta = (random.random() * 2 - 1) * 60  # Da -60 a +60 gradi
        rcs = 0.5 + random.random() * 2
        snr = 15.0 + random.random() * 20
        speed = 0.0

    # Calcolo posizione e riclassificazione finale
    x, y = calculate_position(distance, theta)
    final_task_id = classify_task(distance, rcs, snr)

    return Radar(
        timestamp=timestamp,
        distance=distance,
        theta=theta,
        x=x,
        y=y,
        speed=speed,
        rcs=rcs,
        snr=snr,
        task_id=final_task_id
    )

if __name__ == "__main__":
    """ file_name = "radar_data.csv"

    fieldnames = [
        "timestamp", "distance", "theta", "x", "y", "speed", "rcs", "snr", "task_id"
    ]

    with open(file_name, mode="w", newline="") as csvfile:
        csvfile.write("sep=,\n")
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader() """

    for _ in range(5):
        p = generate_radar_packet()
        print(p)


           # writer.writerow(asdict(p))

    """ print(f"Simulazione completata. Dati salvati in: {os.path.abspath(file_name)}") """