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
 * AlertButton use Button to show some info, 
 * but it can be use not only for button
 * @param {*} Id the Id of Button
 * @param {*} className the temp className of Button
 * @param {*} innerHTML the temp innerHtml of Button
 * @param {*} delay after delay ms the button will flash back 
 */
function AlertButton(Id, className, innerHTML, delay) {
    let OldclassName = document.getElementById(Id).className
    let OldinnerHTML = document.getElementById(Id).innerHTML
    
    document.getElementById(Id).className = className;
    document.getElementById(Id).innerHTML = innerHTML;

    setTimeout("document.getElementById(\"{Id}\").className = \"{OldclassName}\"".format({"Id":Id,"OldclassName":OldclassName}),delay);
    setTimeout("document.getElementById(\"{Id}\").innerHTML = \"{OldinnerHTML}\"".format({"Id":Id,"OldinnerHTML":OldinnerHTML}),delay);
}

/*
String format Space Start
*/
/**
 * Format String
 * @param {*} args 
 */
String.prototype.format = function (args) {
    var result = this;
    if (arguments.length > 0) {
        if (arguments.length == 1 && typeof (args) == "object") {
            for (var key in args) {
                if (args[key] != undefined) {
                    var reg = new RegExp("({" + key + "})", "g");
                    result = result.replace(reg, args[key]);
                }
            }
        }
        else {
            for (var i = 0; i < arguments.length; i++) {
                if (arguments[i] != undefined) {
                    　　　　　　　　　　　　var reg = new RegExp("({)" + i + "(})", "g");
                    result = result.replace(reg, arguments[i]);
                }
            }
        }
    }
    return result;
}