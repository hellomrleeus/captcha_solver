echo $1
source /usr/anaconda3/activate test
cd /root/captcha_solver
python myTest.py --image $1