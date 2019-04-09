BIN=$1
TIMES=$2
for i in `seq 1 $TIMES`; do
  time $BIN install -c akoi.yml
  rm dummy/ubuntu-desktop-iso*
done
