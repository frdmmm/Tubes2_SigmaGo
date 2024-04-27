# Tubes2_SigmaGo

> Tugas Besar kedua mata kuliah IF2211 Strategi Algoritma ITB 2023/2024

## Table of Contents

- [General Info](#general-information)
- [Technologies Used](#technologies-used)
- [Features](#features)
- [Setup](#setup)
- [Usage](#usage)
- [Project Status](#project-status)
<!-- * [License](#license) -->

## General Information

Projek ini bertujuan untuk menemukan solusi permainan Wikirace dengan menggunakan algoritma pencarian BFS atau IDS
## Technologies Used

- Go - version 1.22.2
- npm - version 10.4.0
- Gin - version 1.9.1
- Goquery - version 1.9.1

## Features

- Menampilkan hasil pencarian solusi permainan Wikirace dengan algoritma BFS
- Menampilkan hasil pencarian solusi permainan Wikirace dengan algoritma IDS


## Setup
Clone repository ini
```
git clone https://github.com/frdmmm/Tubes2_SigmaGo
```
FrontEnd:
```
cd src/FrontEnd
npm install -g http-server
```

## Usage
FrontEnd:

Jalankan di powershell (jika di Windows)
```
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass 
http-server -p 7000 -c-0
```

BackEnd:
```
cd src/BackEnd
go run main.go
```

lalu buka website di 
```
localhost/7000
```

## Project Status

Project is: _complete_

## Contact

Dibuat oleh

- Ahmad Farid Mudrika(13522008)
- Devinzen(13522064)
- Muhammad Fuad Nugraha (10023520)


<!-- Optional -->
<!-- ## License -->
<!-- This project is open source and available under the [... License](). -->

<!-- You don't have to include all sections - just the one's relevant to your project -->