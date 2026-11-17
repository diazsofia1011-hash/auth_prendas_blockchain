const express = require('express');
const bodyParser = require('body-parser');
const cors = require('cors');

const app = express();
const PORT = process.env.PORT || 3001;

app.use(cors());
app.use(bodyParser.json());

// Ruta de prueba para verificar que el backend de la fábrica está activo
app.get('/', (req, res) => {
    jsonResponse = {
        empresa: "Fabricante / Marca (Org1)",
        estado: "Activo y conectado al sistema de autenticación"
    };
    res.json(jsonResponse);
});

// Endpoint para simular el registro de una prenda (Aquí despues conectaremos el SDK de Fabric)
app.post('/api/registrar-prenda', (req, res) => {
    const { idPrenda, marcaId, hashSecreto, metadatoHash } = req.body;

    if (!idPrenda || !marcaId) {
        return res.status(400).json({ error: "Faltan datos obligatorios (idPrenda, marcaId)" });
    }

    console.log(`[FÁBRICA] Registrando nueva prenda en la blockchain: ${idPrenda}`);

    // Simulación de respuesta exitosa de la transacción enviada a los peers de Fabric
    res.status(201).json({
        mensaje: "Prenda registrada exitosamente en la red de Fabric por la Marca",
        datosPrenda: {
            id: idPrenda,
            marcaId: marcaId,
            hashSecreto: hashSecreto,
            estado: "EN_FABRICA"
        }
    });
});

app.listen(PORT, () => {
    console.log(`Backend de la Fábrica corriendo en http://localhost:${PORT}`);
});