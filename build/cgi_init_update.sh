#!/bin/sh
xmlfile="initXml.xml"

ScriptSelf=$(cd "$(dirname "$0")"; pwd)
echo "				setup begin ***********************************"
sync
echo 3 > /proc/sys/vm/drop_caches

drgflydir="/drgfly/webomt/web/www/"
drgflyhttpsdir="/drgfly/httpswww/"

#backup webomt packet
dir=$ScriptSelf
webBackpath="/drgfly/webomt/"
WebPacketInfoPath="/drgfly/webomt/WebPacketInfo.txt"
md5path="/drgfly/webomt/web/websum.md5"

webpath=`find  $dir  -name "*WEBOMT*" | head -n 1 | awk '{ print $1 }'`

if [ "$webpath" == ""  ];then
    dir=/drgfly/webomt/web/www/packetUpdate
    webpath=`find  $dir  -name "*WEBOMT*" | head -n 1 | awk '{ print $1 }'`
    filename=`basename $webpath`
else
    filename=`basename $webpath`
fi
echo $webpath

if [ -f "$WebPacketInfoPath" ]
then
	echo "				$WebPacketInfoPath exist"
    res=`grep "PATH" -w $WebPacketInfoPath`
    if [ "$res" == "" ];then
        echo "PATH=">$WebPacketInfoPath
	    echo "CNT=0">>$WebPacketInfoPath
    fi
    res=`grep "CNT" -w $WebPacketInfoPath`
    if [ "$res" == "" ];then
        echo "PATH=">$WebPacketInfoPath
	    echo "CNT=0">>$WebPacketInfoPath
    fi
else
	echo "				touch $WebPacketInfoPath"
	touch $WebPacketInfoPath
	echo "PATH=">$WebPacketInfoPath
	echo "CNT=0">>$WebPacketInfoPath
fi

if [ "$webpath" != "" ];then
    #echo "PATH="$filename>$WebPacketInfoPath
	sed -i '1c PATH='"$filename"'' $WebPacketInfoPath
	rm $webBackpath*.zip
    cp $webpath $webBackpath$filename
fi

sync
sync

#+++++++++++++++
if [ -d "$drgflydir" ]
then
	echo "				$drgflydir exist"
else
	echo "				mkdir $drgflydir"
	mkdir $drgflydir
fi

if [ -d "$drgflyhttpsdir" ]
then
	echo "				$drgflyhttpsdir exist"
else
	echo "				mkdir $drgflyhttpsdir"
	mkdir $drgflyhttpsdir
fi

tempdir="/drgfly/webomt/web/www/UploadFiles"
if [ -d "$tempdir" ]
then
	echo "				$tempdir exist"
else
	echo "				mkdir $tempdir"
	mkdir $tempdir
fi

tempdir="/drgfly/webomt/web/www/LogFiles"
if [ -d "$tempdir" ]
then
	echo "				$tempdir exist"
else
	echo "				mkdir $tempdir"
	mkdir $tempdir
fi

tempdir="/drgfly/webomt/web/www/ConfigFiles"
if [ -d "$tempdir" ]
then
	echo "				$tempdir exist"
else
	echo "				mkdir $tempdir"
	mkdir $tempdir
fi

tempdir="/drgfly/webomt/web/www/Vertion"
if [ -d "$tempdir" ]
then
	echo "				$tempdir exist"
else
	echo "				mkdir $tempdir"
	mkdir $tempdir
fi

#ls /drgfly/webomt/web/www/ | grep -v UploadFiles | grep -v LogFiles | xargs -t -I  {} cp /drgfly/webomt/web/www/{} /drgfly_bak/webomt/web/www/ -r
#+++++++++++ add file protect step 1 begin ++++++++++++++++++++++++++++++
webprotectdir="/drgfly/webomt/lastwebback"
unzipdir=$ScriptSelf
if [ -d "$webprotectdir" ]
then
	echo "				$webprotectdir exist"
else
	echo "				mkdir $webprotectdir"
	mkdir $webprotectdir
fi

rm -rf /drgfly/webomt/lastwebback/*

#add new web
chmod 755 -R $ScriptSelf/*

cp -r $ScriptSelf/webomt/web/www/config/*.ini /drgfly/webomt/web/www/cgi-bin/
sync
#mv $ScriptSelf/webomt/web/www/static  $drgflyhttpsdir
#mv $ScriptSelf/webomt/web/www/v-js    $drgflyhttpsdir
#mv $ScriptSelf/webomt/web/www/config  $drgflyhttpsdir
#mv $ScriptSelf/webomt/web/www/cgi-bin $drgflyhttpsdir
mv $ScriptSelf/webomt/web/www/*	$drgflyhttpsdir
rm -r /drgfly/webomt/bin
mv $ScriptSelf/webomt/bin     /drgfly/webomt/
chmod 755 /drgfly/webomt/bin/*
rm $md5path
mv $ScriptSelf/webomt/web/websum.md5  $md5path
sync
sync
echo "				setup end***********************************************"

exit 0
