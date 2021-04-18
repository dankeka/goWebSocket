var usernameInput: HTMLInputElement = document.getElementById("username") as HTMLInputElement
var password1Input: HTMLInputElement = document.getElementById("password1") as HTMLInputElement
var password2Input: HTMLInputElement = document.getElementById("password2") as HTMLInputElement

var errBlock: HTMLDivElement = document.getElementById("registerErrors") as HTMLDivElement

function submitForm(e: HTMLFormElement): boolean {
  if(password1Input.value !== password2Input.value || password1Input.value === "") {
    errBlock.setAttribute("class", "mb-3")
    errBlock.innerHTML = '<span class="text-danger">Пароли не совпадают!</span>'

    return false;
  } else {
    errBlock.removeAttribute("class")
    errBlock.innerHTML = ""
  }

  return true;
}