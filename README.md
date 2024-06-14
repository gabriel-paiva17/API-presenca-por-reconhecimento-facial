https://documenter.getpostman.com/view/34092137/2sA3QngYxU


# USUARIO:

pode haver mais de 1 com o mesmo nome, mas nao pode haver mais de 1 com email mesmo email

# GRUPO:

nao pode haver grupos com mesmo nome do mesmo usuario

# MEMBRO:

nao pode haver um membro do mesmo usuario, no memsmo grupo, com nome igual

# SESSAO: 

nao podem haver sessoes do mesmo usuario com o mesmo nome no mesmo grupo


# REFLETIR:

- Dados de detalhes do grupo sao extremamente voláteis, já que são calculados com base nos dados de sessoes.
talvez seja uma boa ideia retirar attendance do BD de grupo. e ser adicionado somente na response de GET detalhes (calculado em tempo real)
- Assim, membros teriam attendances nas sessoes, mas a attendance deles no grupo, é calculada somente quando voce busca por detalhes do grupo 

# Não esquecer:

- Criar uma funcao chamada update group, quando ela for chamada, a presença de cada aluno no grupo, sera atualizada para a soma de todas as presenças das sessões (é possível colocar essa função até na função que encontra um grupo especifico)
- Implementar função para visualizar sessoes pendentes (isto é, iniciadas e nao finalizadas)
- Implementar forma de verificar se a face foi validada ou nao durante a sessao (se foi, mostrar imagem capturada e adicionar presenças)

# Observações:

- Visão Computacional: espera-se que a base64 venha de um .jpg (multilinhas)
- Ainda não consegui fazer uma função que verifica se é um .jpg multilinha ou um .jpeg em uma linha
- Função que somente checa autenticação não esta funcionando (apesar de no back estar funcionando, o front nao captura as responses)

