document.getElementById('login-from').addEventListener('submit',function(event){
    event.preventDefault();
    var formData = new formData(this);

    fetch('/login',{
        method:'POST',
        body: formData
    }).then(response =>{
        if (response.ok){
            return response.text();
        }
        throw new Error('Giriş işlemi başarısız');
    }).then(data =>{
        alert(data);
    }).catch(error => {
        alert(error.message);
    });
});

document.getElementById('register-from').addEventListener('submit',function(event){
    event.preventDefault();
    var formData = new formData(this);

    fetch('/register',{
        method:'POST',
        body: formData
    }).then(response =>{
        if (response.ok){
            return response.text();
        }
        throw new Error('Kayıt işlemi başarısız');
    }).then(data =>{
        alert(data);
    }).catch(error => {
        alert(error.message);
    });
});


// main.js dosyasına eklenecek JavaScript kodu
function navigateToCategory() {
    const dropdown = document.getElementById('categoryDropdown');
    const selectedValue = dropdown.value;

    if (selectedValue) {
        window.location.href = selectedValue;
    }
}
