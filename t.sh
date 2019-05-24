echo $1
source /usr/anaconda3/bin/activate test
cd /root/captcha_solver
python myTest.py --image $1