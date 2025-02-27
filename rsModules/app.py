from flask import Flask, request
import function.cpu_stress_module as cpu_stress_
import function.gpu_stress_module as  gpu_stress_
import function.memory_stress_module as memory_stress_
import function.disk_stress_module as disk_stress_
import function.network_stress_module as network_stress_
import function.all_in_one_stress_module as all_in_one_
app = Flask(__name__)

@app.route('/')
def home():
   return 'This is Home!'

#Ex.) http://10.0.2.193:5000/cpuStress?duration=5&cpu_num=2
@app.route('/cpu_stress')
def cpu_stress_function():
    parameter_dict = request.args.to_dict(flat=False)
    # cpu_stress_.cpu_stress_ng_func(parameter_dict.get('cpu_num')[0], int(parameter_dict['duration'][0]), parameter_dict.get('percentage')[0])
    cpu_stress_.run_process(parameter_dict.get('cpu_num')[0], int(parameter_dict['duration'][0]), parameter_dict.get('percentage')[0])
    return "CPU Stress Test Completed", 200

@app.route('/gpu_stress')
def gpu_stress_test_func():
    parameter_dict = request.args.to_dict(flat=False)
    # gpu_stress_.gpu_stress_all(int(parameter_dict['duration'][0]))
    gpu_stress_.run_process(int(parameter_dict['duration'][0]))
    return "GPU Stress Test Completed", 200

#Ex.) http://10.0.2.193:5000/memoryStress?duration=5&mem_amount=600
@app.route('/memory_stress')
def memory_stress_test_func():
    parameter_dict = request.args.to_dict(flat=False)
    memory_stress_.run_process(int(parameter_dict['duration'][0]), int(parameter_dict['mem_amount'][0]))
    return "Memory Stress Test Completed", 200

#Ex.) http://10.0.2.193:5000/networkStress?duration=5&mode=preprocess&net_url=localhost&net_port=5000&network_mode=prep
@app.route('/network_stress')
def network_stress_test_func():
    parameter_dict = request.args.to_dict(flat=False)
    network_stress_.run_process(int(parameter_dict['duration'][0]), parameter_dict['net_url'][0], parameter_dict['net_port'][0], parameter_dict['network_mode'][0])
    return "Network Stress Test Completed", 200

'''
dictionary examples
{
    'cpu_stress': ['True'],
    'gpu_stress': ['False'],
    'memory_stress': ['False'],
    'disk_stress': ['True'],
    'network_stress': ['True'],
    'duration': ['5'],
    'mode': ['preprocess'],
    'net_url': ['localhost'],
    'net_port': ['5000'],
    'network_mode': ['prep'],
    'cpu_num': ['2'],
    'mem_amount': ['600'],
    'size_mb': ['100']
}
'''
@app.route('/aio')
def all_in_one_stress_test_function():
    parameter_dict = request.args.to_dict(flat=False)
    print(parameter_dict)
    print('type of: ', type(parameter_dict.get('duration')))
    all_in_one_.all_in_one_test_func(
        parameter_dict
    )

    return 'all_in_one_stress_test_function'

if __name__ == '__main__':
   app.run('0.0.0.0',port=5000,debug=True)
