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
document.addEventListener("DOMContentLoaded", function () {
    const toggleButtons = document.querySelectorAll(".toggle-comments");
    toggleButtons.forEach(button => {
        button.addEventListener("click", function () {
            const commentsContainer = this.nextElementSibling;
            commentsContainer.style.display = commentsContainer.style.display === "none" ? "block" : "none";
        });
    });
});
// main.js
$(document).ready(function() {
    $('.comment_form').submit(function(event) {
        event.preventDefault();
        
        var form = $(this);
        var url = form.attr('action');
        var formData = form.serialize();

        $.ajax({
            type: 'POST',
            url: url,
            data: formData,
            dataType: 'json', // JSON veri türü bekliyoruz
            success: function(response) {
                // Eski başarı veya hata mesajlarını temizle
                form.find('.message').remove();
                
                // Yorum başarıyla eklenmişse
                if (response.message === "Comment posted successfully") {
                    var successMessage = '<p class="message success_message">' + response.message + '</p>';
                    form.prepend(successMessage); // Formun başına ekle
                    form.find('textarea').val(''); // Yorum alanını temizle
                } else {
                    // Başka bir mesaj döndüyse (genelde bu durumu özel olarak işlemek gerekir)
                    console.error('Unexpected response:', response);
                }
            },
            error: function(xhr, status, error) {
                console.error('AJAX error:', error);
                var errorMessage = '<p class="message error_message">Yorum eklenirken bir hata oluştu</p>';
                form.prepend(errorMessage); // Formun başına ekle
            }
        });
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

