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

if __name__ == '__main__':
   app.run('0.0.0.0',port=5000,debug=True)
