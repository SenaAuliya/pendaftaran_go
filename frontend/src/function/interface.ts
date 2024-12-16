export interface Jurusan {
    ID: number;
    Nama: string;
    KodeJurusan: string;
}

export interface AsalSekolah {
    ID: number;
    Nama: string;
}

export interface TahunAjaran {
    ID: number;
    Tahun: string;
}

export interface Siswa {
    ID: number;
    Nama: string;
    NIK: number;
    NISN: number;
    AsalSekolahID: number; // Changed to number
    Nilai: number;
    JarakTempuh: number;
    Alamat: string;
    TempatLahir: string;
    TanggalLahir: string;
    Gender: string;
    FotoURL: string;
}

export interface UserData {
    ID: number;
    Username: string;
    Password: string;
}


export interface User {
    Username: string;
    Role?: string;
    // Add other properties of the User object if necessary
}

export interface Guru {
    ID: number;
    Nama: string;
    KodeGuru: string;
}

export interface Pendaftaran {
    ID: number;
    GuruID: number;         // Changed to number
    SiswaID: number;        // Changed to number
    TahunID: number;        // Changed to number
    JurusanID: number;      // Changed to number
    TanggalPendaftaran: string; // Assuming this remains as string
    TipePendaftaran: string;
    limit: number
}

export interface ApiResponse<T> {
    data: T[];
    totalPages: number;
    currentPage: number;
    totalItems: number;
}
