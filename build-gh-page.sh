ng build --base-href /harpa/
rm -R docs
mv dist/harpa docs
cp src/youtube.html docs/