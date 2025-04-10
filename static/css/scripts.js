function deletePost(postID) {
    fetch(`/posts/${postID}`, {
        method: "DELETE"
    }).then(response => {
        if (response.ok) {
            alert("Пост удален");
            location.reload(); 
        } else {
            alert("Ошибка при удалении поста");
        }
    });
}

function previewImage(event) {
    const reader = new FileReader();
    reader.onload = function() {
        const avatar = document.getElementById('avatar');
        const avatarText = document.getElementById('avatar-text');

        avatar.src = reader.result; 
        avatarText.classList.add('hidden');  
    };
    reader.readAsDataURL(event.target.files[0]);
}

