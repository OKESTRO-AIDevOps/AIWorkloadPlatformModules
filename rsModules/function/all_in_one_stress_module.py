from multiprocessing import Process, active_children, cpu_count, Pipe
import os
import sys
import time
import signal

import function.cpu_stress_module as cpu_stress_
import function.memory_stress_module as memory_stress_
import function.disk_stress_module as disk_stress_
import function.network_stress_module as network_stress_


def all_in_one_test_func(parameter_dict):
    mp_list = []
    time = int(parameter_dict.get('duration')[0])
    if 'cpu_stress' in parameter_dict:
        if parameter_dict.get('cpu_stress')[0] == 'True':
            try:
                mp_cpu_stress = Process(
                target=cpu_stress_.cpu_stress_ng_func,
                args=(
                    parameter_dict.get('cpu_num')[0],
                    time,
                    parameter_dict.get('percentage')[0]
                    )
                )
                mp_list.append(mp_cpu_stress)
            except ValueError:
                print('[Error] Input parameter value error')
                pass
        else:
            pass
    else:
        pass

    if 'gpu_stress' in parameter_dict:
        pass

    if 'memory_stress' in parameter_dict:
        if parameter_dict.get('memory_stress')[0]  == 'True':
            try:
                mp_memory_stress = Process(
                target=memory_stress_.memory_stress_func,
                args=(
                    time,
                    int(parameter_dict.get('mem_amount')[0])
                    )
                )
                mp_list.append(mp_memory_stress)
            except ValueError:
                print('[Error] Input parameter value error')
                pass
        else:
            pass
    else:
        pass

    if 'disk_stress' in parameter_dict:
        if parameter_dict.get('disk_stress')[0]  == 'True':
            try:
                mp_disk_stress = Process(
                target=disk_stress_.disk_stress_func,
                args=(
                    time,
                    int(parameter_dict.get('size_mb')[0])
                    )
                )
                mp_list.append(mp_disk_stress)
            except ValueError:
                print('[Error] Input parameter value error')
                pass
        else:
            pass
    else:
        pass
    if 'network_stress' in parameter_dict:
        if parameter_dict.get('network_stress')[0] == 'True':
            try:
                mp_network_stress = Process(
                target=network_stress_.network_stress_func,
                args=(
                    time,
                    parameter_dict.get('net_url'),
                    parameter_dict.get('net_port'),
                    parameter_dict.get('network_mode')
                    )
                )
                mp_list.append(mp_network_stress)
            except ValueError:
                print('[Error] Input parameter value error')
                pass
        else:
            pass
    else:
        pass
    for mp in mp_list:
        mp.start()

    return 'all_in_one_test_func'
