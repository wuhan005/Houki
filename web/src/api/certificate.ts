import axios from 'axios';

export interface CertificateMetadata {
  issuer: string;
  validFrom: Date;
  validTo: Date;
  publicKeyAlgorithm: string;
  serialNumber: string;
  signatureAlgorithm: string;
}

export interface GetCertificateResp {
  certificate: string;
  privateKey: string;
  metadata: CertificateMetadata;
}

export function getCertificate() {
  return axios.get<GetCertificateResp, GetCertificateResp>('/api/certificate');
}

export interface UpdateCertificateData {
  certificate: string;
  privateKey: string;
}

export function updateCertificate(data: UpdateCertificateData) {
  return axios.put<string, string>('/api/certificate', data);
}
