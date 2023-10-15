#VER=3.0.1-beta.2310151130
#EXT=.rpm

URL=https://github.com/nuvolaris/nuv/releases/download/$VER/nuv_${VER}_${ARCH}${EXT} >nuv${EXT}
echo $URL
case "$EXT" in 
  *.deb)
    apt-get update && apt-get install -y curl
    export INST="dpkg -i"
  ;;
  *.rpm) 
    export INST="rpm -i"
  ;;
esac

curl -sL "$URL" >nuv${EXT}
$INST "nuv${EXT}"

nuv -v 

