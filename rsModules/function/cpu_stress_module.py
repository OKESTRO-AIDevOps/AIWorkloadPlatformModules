from function.pystress import pystress
import subprocess

def cpu_stress_func():
    # print(f"\n cpu cores using cpu stress test : {args.cpu_num}")
    # print(f"Generating CPU load. (Duration: {duration} s)")
    # pystress(exec_time=duration,proc_num=args.cpu_num) # pystress CPU 부하 테스트
    print(f"\n cpu cores using cpu stress test : {2}")
    print(f"Generating CPU load. (Duration: {10} s)")
    pystress(exec_time=10, proc_num=2)
    print(f"CPU Stress Process created. (Duration: {10} m)")
    return 'cpu stress test'

def cpu_stress_func2(cpu_num:int, duration: int):
    pystress(exec_time=duration, proc_num=cpu_num)
    return 'url data parse test - cpu stress test'

def cpu_stress_ng_func(cpu_num:str, duration:int, percentage: str):
    command = ["stress-ng", "--cpu", str(cpu_num), "--cpu-load", str(percentage), "--timeout", str(duration)]
    subprocess.run(command)%
