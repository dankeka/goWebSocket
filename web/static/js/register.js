var usernameInput = document.getElementById("username");
var password1Input = document.getElementById("password1");
var password2Input = document.getElementById("password2");
var errBlock = document.getElementById("registerErrors");
function submitForm(e) {
    if (password1Input.value !== password2Input.value || password1Input.value === "") {
        errBlock.setAttribute("class", "mb-3");
        errBlock.innerHTML = '<span class="text-danger">Пароли не совпадают!</span>';
        return false;
    }
    else {
        errBlock.removeAttribute("class");
        errBlock.innerHTML = "";
    }
    return true;
}
