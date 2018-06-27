##author      : anil.khadwal@gmail.com 
##description : to check available open connections over given range of IP_ADDRESS

SERVER=127.0.0  
PORT=23001

#i=1
for ((i=0 ;i<=255; i++))
do


`nc -z -v -w5 $SERVER.$i $PORT`
result1=$?

#echo 'checking on '$i

if [  "$result1" != 0 ]; then
  echo  
else
  echo  $SERVER.$i:$PORT >> available_con.txt 
fi

done
