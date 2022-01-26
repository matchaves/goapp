# goapp


192.168.50.32:32001/ PrintEnv
printa as varias de ambiente FOO no logs


192.168.50.32:32001/post Postparams
{"name": "mateus", 
 "club": "bumba",
 "age": 8
}			

192.168.50.32:32001/param1/{param1}/param2/{param2} Urlparams
retorna como resposta os params que vc passa na url

192.168.50.32:32001/to/{app}  MakeRequest

//export URLAPP2=http://172.18.0.6/post
//export URLAPP2=http://172.18.0.6/auth

variaveis das app
IP_APP1 = 192.xxx.xx.xx
IP_APP2
IP_AUTH
