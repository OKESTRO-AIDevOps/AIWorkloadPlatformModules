import time
from multiprocessing import Process

def disk_stress_func(duration:int, size_mb:int=500):
    print(f"\n amount of disk stress : {size_mb}")
    print(f"Generating disk I/O load. (Duration: {duration} s)")
    #file write
    with open('tmp','wb') as f:
        f.write(b'\0' * (size_mb * 1024 * 1024))
        time.sleep(duration//2)
    #file read
    with open('tmp','rb') as f:
        data= f.read()
        time.sleep(duration//2)

def multi_processing_disk_stress_func(duration:int, size_mb:int=500):
    p = Process(target=disk_stress_func, args=(duration, size_mb))
    p.start()