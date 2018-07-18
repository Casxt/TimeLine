/**
 * Cookies Addons
 * TODO: split("; ") should have change and not test in Edge and IOS
 * Work with Edge
 */
Object.defineProperty(window, "Cookies", {
    get: function () {
        return document.cookie.split("; ").reduce(function (cookies, cookie) {
            const cookieitem = cookie.split("=");
            cookies[cookieitem[0]] = unescape(cookieitem[1]);
            return cookies
        }, {});
    }
});

/**
 * FormData Addons
 */
FormData.prototype.urlencode = function () {
    var params = new URLSearchParams();
    for (var pair of this.entries()) {
        typeof pair[1] == 'string' && params.append(pair[0], pair[1]);
    }
    return params.toString();
}
/**
 * ToArray Trans A FormData to Array
 * @param {*} args 
 */
FormData.prototype.ToArray = function (args) {
    let Data = {};
    this.forEach((value, key) => Data[key] = value);
    return Data;
}

async function JsonRequest(httpmethod, url, data) {
    const request = new Request(url, {
        credentials: 'include', //接受response中的cookie
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        method: httpmethod,
        body: JSON.stringify(data),
    });

    const response = await fetch(request);
    return await response.json();
}

async function FileRequest(httpmethod, url, data) {
    const request = new Request(url, {
        credentials: 'include', //接受response中的cookie
        headers: {
            'Accept': 'application/json',
            //'Content-Type': 'multipart/form-data'
        },
        method: httpmethod,
        body: data,
    });
    const response = await fetch(request);
    return await response.json();
}


/*
String format Space Start
*/
/**
 * Format String
 * @param {*} args 
 */
String.prototype.format = function (args) {
    
    let result = this;
    if (arguments.length > 0) {
        if (arguments.length == 1 && typeof (args) == "object") {
            for (const key in args) {
                if (args[key] != undefined) {
                    result = result.replace(RegExp("{" + key + "}","g"), args[key]);
                }
            }
        } else {
            for (let i  = 0; i < arguments.length; i++) {
                if (arguments[i] != undefined) {　　　　　　　　　　　　
                    const reg = new RegExp("({)" + i + "(})", "g");
                    result = result.replace("{" + i + "}", arguments[i]);
                }
            }
        }
    }
    return result;
}

function FormtTime(TimeStr){
    return TimeStr.replace("T"," ").replace("Z","")
}

class AnimeButton {
    constructor(buttonId) {
        if (typeof buttonId == "string") {
            this.buttonId = buttonId;
            this.button = document.getElementById(buttonId);
            this.jqButton = $("#" + buttonId);
        }else if (Object.prototype.toString.call(buttonId) == "[object HTMLButtonElement]"){
            this.buttonId = null;
            this.button = buttonId;
            this.jqButton = $(buttonId);
        }
    }
    /**
     * Alert use Button to show some info, 
     * but it can be use not only for button
     * @param {*} Id the Id of Button
     * @param {*} className the temp className of Button
     * @param {*} innerHTML the temp innerHtml of Button
     * @param {*} delay after delay ms the button will flash back 
     */

    OnLoding(className, innerHTML) {
        this.jqButton.addClass(className);
        const OldHtml = this.button.innerHTML
        this.button.disabled = true;
        this.button.innerHTML = innerHTML;
        return () => {
            this.button.disabled = false;
            this.jqButton.removeClass(className);
            this.button.innerHTML = OldHtml;
        }
    }

    Alert(className, innerHTML, delay) {
        const OldclassName = this.button.className
        const OldinnerHTML = this.button.innerHTML

        this.button.className = className;
        this.button.innerHTML = innerHTML;

        setTimeout(() => {
            this.button.className = OldclassName;
            this.button.innerHTML = OldinnerHTML;
        }, delay);
    }
}