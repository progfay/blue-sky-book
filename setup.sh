rm -rf texts
mkdir texts

if [ ! -d "aozorabunko_text" ]; then
  git clone --depth 1 https://github.com/aozorahack/aozorabunko_text.git
fi

for txt_file in $(find "`pwd`" -name *.txt)
do
  f="${txt_file##${PWD}\/aozorabunko_text\/cards/}"
  f="texts/${f//\//__}"
  iconv -c -f sjis -t utf8 "${txt_file}" > "$f" &
done

wait
