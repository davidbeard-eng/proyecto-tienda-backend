import { useState, useEffect } from 'react';

// Títulos legibles para cada una de las 10 consultas solicitadas
const titulosReportes: { [key: string]: string } = {
  r1: "1. Producto más vendido por mes el 2021",
  r2: "2. Producto más económico por tienda",
  r3: "3. Ventas por mes, separadas entre Boletas y Facturas",
  r4: "4. Empleado que ganó más por tienda en 2020 (con Comuna y Cargo)",
  r5: "5. La tienda que tiene menos empleados",
  r6: "6. El vendedor con más ventas por mes",
  r7: "7. El vendedor que ha recaudado más dinero para la tienda por año",
  r8: "8. El vendedor con más productos vendidos por tienda",
  r9: "9. El empleado con mayor sueldo por mes",
  r10: "10. La tienda con menor recaudación por mes"
};

type ReportData = { [key: string]: any }[];

export default function App() {
  const [reporteActivo, setReporteActivo] = rangeReporte("r1");
  const [datos, setDatos] = useState<ReportData>([]);
  const [error, setError] = useState<string | null>(null);
  const [cargando, setCargando] = useState<boolean>(false);

  // Efecto que se dispara cada vez que cambias de reporte en el menú
  useEffect(() => {
    async function cargarDatos() {
      setCargando(true);
      setError(null);
      try {
        const response = await fetch(`http://localhost:8080/api/reportes/${reporteActivo}`);
        if (!response.ok) throw new Error("Error en la respuesta del servidor");
        const json = await response.json();
        setDatos(json || []);
      } catch (err: any) {
        setError("No se pudo conectar con el backend de Go. Asegúrate de que esté corriendo.");
        setDatos([]);
      } finally {
        setCargando(false);
      }
    }
    cargarDatos();
  }, [reporteActivo]);

  // Helper para alternar entre reportes de forma segura
  function rangeReporte(initialValue: string): [string, (val: string) => void] {
    const [val, setVal] = useState(initialValue);
    return [val, setVal];
  }

  // Extraer dinámicamente las columnas del JSON recibido para armar la cabecera de la tabla
  const columnas = datos.length > 0 ? Object.keys(datos[0]) : [];

  return (
    <div style={styles.container}>
      {/* BARRA LATERAL DE NAVEGACIÓN */}
      <aside style={styles.sidebar}>
        <h2 style={styles.logo}>🛒 TBD Panel</h2>
        <p style={styles.subtext}>Control 1 - Grupo 1</p>
        <nav style={styles.nav}>
          {Object.keys(titulosReportes).map((key) => (
            <button
              key={key}
              onClick={() => setReporteActivo(key)}
              style={{
                ...styles.navButton,
                backgroundColor: reporteActivo === key ? '#2563eb' : 'transparent',
                color: reporteActivo === key ? '#ffffff' : '#94a3b8'
              }}
            >
              Consulta {key.toUpperCase().replace('R', '')}
            </button>
          ))}
        </nav>
      </aside>

      {/* ÁREA PRINCIPAL DE CONTENIDO */}
      <main style={styles.mainContent}>
        <header style={styles.header}>
          <h1 style={styles.title}>{titulosReportes[reporteActivo]}</h1>
        </header>

        <section style={styles.contentCard}>
          {cargando && <p style={styles.infoMessage}>🔄 Cargando datos desde la API de Go...</p>}
          {error && <p style={styles.errorMessage}>⚠️ {error}</p>}
          
          {!cargando && !error && datos.length === 0 && (
            <p style={styles.infoMessage}>No hay datos registrados para este reporte en la base de datos.</p>
          )}

          {/* TABLA DINÁMICA */}
          {!cargando && datos.length > 0 && (
            <div style={styles.tableWrapper}>
              <table style={styles.table}>
                <thead>
                  <tr style={styles.thRow}>
                    {columnas.map((col) => (
                      <th key={col} style={styles.th}>
                        {col.toUpperCase().replace('_', ' ')}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {datos.map((fila, index) => (
                    <tr key={index} style={index % 2 === 0 ? styles.trEven : styles.trOdd}>
                      {columnas.map((col) => (
                        <td key={col} style={styles.td}>
                          {typeof fila[col] === 'number' && col.includes('recaudado') || col.includes('precio') || col.includes('sueldo') || col.includes('ganado')
                            ? `$${fila[col].toLocaleString('es-CL')}` 
                            : fila[col]}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </main>
    </div>
  );
}

// --- ESTILOS EN LÍNEA ESTÁNDAR (Para no depender de CSS externo) ---
const styles = {
  container: { display: 'flex', height: '100vh', backgroundColor: '#0f172a', fontFamily: 'system-ui, sans-serif', color: '#f8fafc' },
  sidebar: { width: '260px', backgroundColor: '#1e293b', padding: '24px', display: 'flex', flexDirection: 'column' as const, borderRight: '1px solid #334155' },
  logo: { fontSize: '20px', fontWeight: 'bold', margin: 0, color: '#38bdf8' },
  subtext: { fontSize: '12px', color: '#64748b', marginTop: '4px', marginBottom: '24px' },
  nav: { display: 'flex', flexDirection: 'column' as const, gap: '8px', flex: 1, overflowY: 'auto' as const },
  navButton: { padding: '10px 14px', border: 'none', borderRadius: '6px', textAlign: 'left' as const, cursor: 'pointer', fontSize: '14px', fontWeight: 500, transition: 'all 0.2s' },
  mainContent: { flex: 1, padding: '40px', display: 'flex', flexDirection: 'column' as const, overflowY: 'auto' as const },
  header: { marginBottom: '24px' },
  title: { fontSize: '24px', fontWeight: 'bold', margin: 0, color: '#f1f5f9' },
  contentCard: { backgroundColor: '#1e293b', borderRadius: '12px', padding: '24px', border: '1px solid #334155', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' },
  tableWrapper: { overflowX: 'auto' as const },
  table: { width: '100%', borderCollapse: 'collapse' as const, textAlign: 'left' as const, fontSize: '14px' },
  thRow: { backgroundColor: '#334155' },
  th: { padding: '12px 16px', fontWeight: 600, color: '#38bdf8', textTransform: 'uppercase' as const, letterSpacing: '0.05em' },
  td: { padding: '12px 16px', borderBottom: '1px solid #334155', color: '#cbd5e1' },
  trEven: { backgroundColor: '#1e293b' },
  trOdd: { backgroundColor: '#1e293b', opacity: 0.9 },
  errorMessage: { color: '#ef4444', backgroundColor: '#451a03', padding: '12px', borderRadius: '6px', margin: 0 },
  infoMessage: { color: '#94a3b8', margin: 0, textAlign: 'center' as const, padding: '20px' }
};