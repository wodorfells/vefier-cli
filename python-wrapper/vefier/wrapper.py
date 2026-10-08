import os
import sys
import platform
import subprocess
import urllib.request
import zipfile
import tarfile
import shutil

VERSION = "1.0.0"
REPO = "vefier/vefier-cli"

def get_release_url():
    sys_name = platform.system().lower()
    arch = platform.machine().lower()
    
    if sys_name == "darwin": os_name = "darwin"
    elif sys_name == "windows": os_name = "windows"
    else: os_name = "linux"
        
    if arch in ["x86_64", "amd64"]: arch_name = "x86_64"
    elif arch in ["arm64", "aarch64"]: arch_name = "arm64"
    else: arch_name = arch
        
    ext = "zip" if os_name == "windows" else "tar.gz"
    # Для теста URL будет ложным, пока нет реального репозитория
    return f"https://github.com/{REPO}/releases/download/v{VERSION}/vefier_{os_name}_{arch_name}.{ext}", ext

def download_and_extract(url, ext, bin_dir, exe_name):
    print(f"[INFO] Загрузка VeFier v{VERSION} из GitHub Releases...")
    archive_path = os.path.join(bin_dir, f"download.{ext}")
    
    try:
        urllib.request.urlretrieve(url, archive_path)
    except Exception as e:
        print(f"[ERROR] Ошибка загрузки ({url}): {e}")
        print("Убедитесь, что релиз существует на GitHub.")
        sys.exit(1)
        
    print("[INFO] Распаковка...")
    if ext == "zip":
        with zipfile.ZipFile(archive_path, 'r') as zf:
            zf.extractall(bin_dir)
    else:
        with tarfile.open(archive_path, 'r:gz') as tf:
            tf.extractall(bin_dir)
            
    os.remove(archive_path)
    if platform.system() != "Windows":
        os.chmod(os.path.join(bin_dir, exe_name), 0o755)

def main():
    # Сохраняем бинарник в домашнюю папку (кэш)
    bin_dir = os.path.expanduser("~/.vefier/bin")
    os.makedirs(bin_dir, exist_ok=True)
    
    exe_name = "vefier.exe" if platform.system() == "Windows" else "vefier"
    bin_path = os.path.join(bin_dir, exe_name)
    
    # Если бинарника нет, скачиваем его
    if not os.path.exists(bin_path):
        url, ext = get_release_url()
        download_and_extract(url, ext, bin_dir, exe_name)
    
    # Передаем управление Go-приложению
    try:
        sys.exit(subprocess.call([bin_path] + sys.argv[1:]))
    except KeyboardInterrupt:
        sys.exit(0)

if __name__ == "__main__":
    main()