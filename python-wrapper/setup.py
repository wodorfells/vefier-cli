from setuptools import setup, find_packages

setup(
    name="vefier-cli",
    version="1.0.0",
    description="Универсальный AI-агент для терминала (TUI)",
    long_description="Обертка для загрузки скомпилированного Go-бинарника VeFier CLI.",
    author="VeFier Studio",
    packages=find_packages(),
    entry_points={
        "console_scripts": [
            "vefier=vefier.wrapper:main",
        ],
    },
    install_requires=[],
)