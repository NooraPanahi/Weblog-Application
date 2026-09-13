const themeToggle = document.getElementById("themeToggle");
const savedTheme = localStorage.getItem("theme");

if (savedTheme === "light") 
    document.documentElement.classList.remove("dark");
 else 
    document.documentElement.classList.add("dark");

themeToggle.textContent = document.documentElement.classList.contains("dark")? "☀️" : "🌙";

themeToggle.addEventListener("click", () => {

    document.documentElement.classList.toggle("dark");
    const isDark = document.documentElement.classList.contains("dark");
    themeToggle.textContent = isDark? "☀️" : "🌙";

    localStorage.setItem(
        "theme", isDark ? "dark" : "light"
    );

});