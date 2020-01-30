(function() {
    const loginView = document.getElementById("login");
    const loginLink = document.getElementById("login-link");
    const loginForm = document.getElementById("login-form");
    const loginError = document.getElementById("login-errors");

    const registerView = document.getElementById("register");
    const registerLink = document.getElementById("register-link");
    const registerForm = document.getElementById("register-form");
    const registerError = document.getElementById("register-errors");

    const overrideSubmit = fetch !== undefined;

    var activeView = null;

    loginView.className = "form-hidden";
    registerView.className = "form-hidden";

    function switchView(newView) {
        if (activeView == newView) {
            return;
        }
        function fadeOut() {
            if (activeView != null) {
                activeView.className = "form-hiding";
                setTimeout(fadeIn, 500);
            } else {
                fadeIn();
            }
        }
        function fadeIn() {
            if (activeView != null) {
                activeView.className = "form-hidden";
            }
            activeView = newView;
            if (activeView != null) {
                activeView.className = "form-showing";
                setTimeout(fadedIn, 500);
            }
        }
        function fadedIn() {
            activeView.className = "form-show";
        }

        fadeOut();
    }

    function setError(element, text) {
        element.innerHTML = "";
        if (text != "") {
            let alert = document.createElement("div");
            alert.className = "callout alert";
            alert.setAttribute("data-alert", "");
            alert.innerText = text;
            element.appendChild(alert);
        }
    }

    switch (window.location.hash) {
    case "#register":
        switchView(registerView);
        break;
    default:
    case "#login":
        switchView(loginView);
        break;
    }

    loginLink.addEventListener("click", function(evt) {
        switchView(loginView);
    });

    registerLink.addEventListener("click", function(evt) {
        switchView(registerView);
    });
})();