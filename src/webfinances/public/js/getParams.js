
function getParams() {
  let scripts = document.getElementsByTagName('script');
  let script = scripts[scripts.length - 1];
  let queryString = script.src.replace(/^[^\?]+\??/, '');
  let Params = new Object();
  if (!queryString) {
    return Params;  //Return an empty object.
  }
  let Pairs = queryString.split(/[;&]/);
  for (let i = 0; i < Pairs.length; ++i) {
    let KeyVal = Pairs[i].split('=');
    if (!KeyVal || KeyVal.length != 2) {
      continue;
    }
    let key = decodeURI(KeyVal[0]);
    let val = decodeURI(KeyVal[1]);
    val = val.replace(/\+/g, ' ');
    //Change 'my-string' to my-string.
    Params[key] = val.substring(1, val.length - 1);
  }
  return Params;
}
