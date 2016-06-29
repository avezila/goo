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

    // Check if the XMLHttpRequest object has a "withCredentials" property.
    // "withCredentials" only exists on XMLHTTPRequest2 objects.
    xhr.open(method, url, true);

  } else if (typeof XDomainRequest != "undefined") {

    // Otherwise, check if XDomainRequest.
    // XDomainRequest only exists in IE, and is IE's way of making CORS requests.
    xhr = new XDomainRequest();
    xhr.open(method, url);

  } else {

    // Otherwise, CORS is not supported by the browser.
    xhr = null;

  }
  return xhr;
}

const SERV = process.env.SERVER;
export function sendUrlForShort(url) {
    if (!url.match(/^\w+\:\//) && !url.match(/^http/))
        url = "http://"+url;
    return function(dispatch) {
        var xhr = createCORSRequest("POST", SERV+"/putUrl");
        xhr.onload = function() {
            if (xhr.status === 200) {
                dispatch(receiveShortenUrl(url,SERV+"/"+xhr.responseText));
            }
            else {
                console.log(xhr.status,xhr.statusText)
            }
        };
        xhr.send(url);
        return;
    }
}
