import * as types from '../constants/ActionTypes';


export function receiveShortenUrl(long,short){
    return {
        type: types.RECEIVE_SHORTEN_URL,
        long,
        short
    };
}

function createCORSRequest(method, url) {
  var xhr = new XMLHttpRequest();
  if ("withCredentials" in xhr) {
    xhr.open(method, url, true);
  } else if (typeof XDomainRequest != "undefined") {
    xhr = new XDomainRequest();
    xhr.open(method, url);
  } else {
    xhr = null;
  }
  return xhr;
}

const SERVER = process.env.SERVER;

export function sendUrlForShort(url) {
  if (!url.match(/^\w+\:\//) && !url.match(/^http/))
    url = "http://"+url;
  return function(dispatch) {
    var xhr = createCORSRequest("POST", window.location.href+"putUrl");
    xhr.onload = function() {
      if (xhr.status === 200) {
          dispatch(receiveShortenUrl(url,window.location.href+xhr.responseText));
      }
    };
    xhr.send(url);
    return;
  }
}
