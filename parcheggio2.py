
import math
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
    task_id: int = 0

def generate_measures() -> dict:
    timestamp = int(time.time() * 1000) # Millisecondi
    if random.random() < P_MAX: # no target
        values = {
            "timestamp": timestamp,
            "range": -1.0,
            "theta": 0.0, 
            "speed": 0.0,
            "rcs": random.uniform(0.0, 0.09),
            "snr": random.uniform(0.0, 9.99)
        }
    else: # oggetto statico
        values = {
            "timestamp": timestamp,
            "range": random.uniform(0.5, RADAR_MAX_RANGE),
            "theta": random.uniform(-60, 60), 
            "speed": 0.0,
            "rcs": random.uniform(0.1, 10.0),
            "snr": random.uniform(10.0, 30.0)
        }

    x, y = calculate_position(values["range"], values["theta"])
    values["x"] = x
    values["y"] = y

    return values

def calculate_position(range: float, theta: float) -> tuple[float, float]:
    """Converte coordinate polari in cartesiane."""
    if range == -1.0:
        return 0.0, 0.0

    # Conversione in radianti per le funzioni trigonometriche
    radians = math.radians(theta)
    x_pos = range * math.sin(radians)
    y_pos = range * math.cos(radians)
    return x_pos, y_pos

if __name__ == "__main__":
    for _ in range(10):
        print(generate_measures())