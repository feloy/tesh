test -f ./path/to/file && echo "file exists" || echo "file does not exist"
if test -f ./path/to/file
then 
  echo file exists
else
  echo file does not exist
fi
