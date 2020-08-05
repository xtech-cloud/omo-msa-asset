go install omo.msa.asset
mkdir _build
mkdir _build/bin

cp -rf /root/go/bin/omo.msa.asset _build/bin/
cp -rf conf _build/
cd _build
tar -zcf msa.asset.tar.gz ./*
mv msa.asset.tar.gz ../
cd ../
rm -rf _build
