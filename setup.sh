if [ ! -d "aozorabunko_text" ]; then
  git clone --depth 1 https://github.com/aozorahack/aozorabunko_text.git
fi

cd aozorabunko_text
git pull
cd ..

rm -rf texts
mkdir texts

for txt_file in $(find "`pwd`/aozorabunko_text" -name "*.txt")
do
  f="${txt_file##${PWD}\/aozorabunko_text\/cards/}"
  f="texts/${f//\//__}"
  iconv -c -f sjis -t utf8 "${txt_file}" > "$f" &
done

wait
