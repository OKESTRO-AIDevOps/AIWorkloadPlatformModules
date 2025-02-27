import subprocess

def memory_stress_func(duration:int, mem_amount:str = "2000"):
    print(f"\n memory amount using memory stress test : {mem_amount}")
    # mem_amount=mem_amount
    print(f"Generating memory load. (Duration: {duration} s)")
    #memaount = test memory size
    #duration: test run time
    command = ["stressapptest", "-s", str(duration), "-M", str(mem_amount)]
    subprocess.run(command)
    print(f"Memory Stress test done. (Duration: {duration} s)")%
