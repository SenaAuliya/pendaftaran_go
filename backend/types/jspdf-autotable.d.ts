declare module 'jspdf-autotable' {
    import { jsPDF } from 'jspdf';
  
    interface AutoTableColumn {
      header: string;
      dataKey: string;
    }
  
    interface AutoTableOptions {
      head?: AutoTableColumn[];
      body?: any[][];
      startY?: number;
      margin?: { top?: number; bottom?: number; left?: number; right?: number };
    }
  
    export function autoTable(doc: jsPDF, options: AutoTableOptions): void;
  }
  