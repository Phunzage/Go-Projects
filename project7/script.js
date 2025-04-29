function getRosePrice() {
    const date = document.getElementById('date').value;
    const resultDiv = document.getElementById('result')
    const loadingDiv = document.getElementById('loading')

// 验证日期
if (!date) {
    alert("请先选择日期");
    return;
}
// 显示加载提示
loadingDiv.style.display = "block";
resultDiv.innerHTML =  '';

// 发起请求
fetch(`/api/rose-price?date=${date}`)
    .then(response => response.json())
    .then(data => {
        loadingDiv.style.display = "none";
        if (data.success) {
            resultDiv.innerHTML = `玫瑰的价格是：￥${data.price}`;
        } else {
            resultDiv.innerHTML = `没有找到该日期的玫瑰价格`;
        }
    })
    .catch(error => {
        loadingDiv.style.display = "none";
        console.error("请求出错：",error);
        alert("请求失败，请稍后重试");
    });
}