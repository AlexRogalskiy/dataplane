import { useState, useMemo, useEffect } from 'react';
import { Box, Grid, Typography } from '@mui/material';
import Search from '../../components/Search';
import { useTable, useGlobalFilter } from 'react-table';
import CustomChip from '../../components/CustomChip';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDocker } from '@fortawesome/free-brands-svg-icons';
import { faSearch } from '@fortawesome/free-solid-svg-icons';
import { formatDate } from '../../utils/formatDate';
// import { useSnackbar } from 'notistack';
import { useGlobalEnvironmentState } from '../../components/EnviromentDropdown';
import WorkerDetailCPU from './WorkerDetailCPU';
import WorkerDetailMemory from './WorkerDetailMemory';
const tableWidth = '1140px';

export default function WorkerDetail() {
    // const { enqueueSnackbar } = useSnackbar();

    // Global environment state with hookstate
    const Environment = useGlobalEnvironmentState();

    // Users state
    const [data, setData] = useState(dummyData);

    // Get workers on load
    // const getUsers = useGetUsers();
    useEffect(() => {
        // retrieveUsers();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Get users
    // const retrieveUsers = async () => {
    //     let users = await getUsers();
    //     !users.errors ? setData(users) : enqueueSnackbar('Unable to retrieve users', { variant: 'error' });
    // };

    const columns = useMemo(
        () => [
            {
                Header: 'Worker',
                accessor: (row) => [
                    row.worker.id,
                    row.worker.status,
                    formatDate(row.worker.lastUpdate),
                    row.worker.queue,
                    row.worker.running,
                    row.worker.succeeded,
                    row.worker.failed,
                ],
                Cell: (row) => <CustomWorker row={row} />,
            },
            {
                Header: 'CPU',
                accessor: (row) => [row.cpu.percentage.toFixed(1), row.cpu.load],
                Cell: (row) => <WorkerDetailCPU row={row} />,
            },
            {
                Header: 'Memoery',
                accessor: (row) => [row.memory.percentage.toFixed(1), row.memory.mb],
                Cell: (row) => <WorkerDetailMemory row={row} />,
            },
        ],
        []
    );

    // Use the state and functions returned from useTable to build your UI
    const { getTableProps, getTableBodyProps, rows, prepareRow, state, setGlobalFilter } = useTable(
        {
            columns,
            data,
        },
        useGlobalFilter
    );

    const { globalFilter } = state;

    useEffect(() => {
        console.log(globalFilter);
    }, [globalFilter]);

    return (
        <Box className="page">
            <Typography component="h2" variant="h2" color="text.primary">
                Workers
            </Typography>

            <Typography variant="subtitle2" mt=".20rem">
                Environment: {Environment.name.get()}
            </Typography>

            <Box mt={4} sx={{ width: tableWidth }}>
                <Grid container mt={4} direction="row" alignItems="center" justifyContent="flex-start">
                    <Grid item display="flex" alignItems="center" sx={{ alignSelf: 'center' }}>
                        <CustomChip amount={rows.length} label="Workers" margin={2} customColor="orange" />
                    </Grid>

                    <Grid item display="flex" alignItems="center" sx={{ marginLeft: 'auto', marginRight: '2px' }}>
                        <FontAwesomeIcon icon={faSearch} style={{ marginRight: 10 }} color="#0000006B" size="xs" />

                        <Search placeholder="Find workers" onChange={setGlobalFilter} width="290px" />
                    </Grid>
                </Grid>
            </Box>

            <Box mt={4} sx={{ width: tableWidth }}>
                <Grid container mt={4} direction="row" alignItems="center" justifyContent="flex-start">
                    <Grid item display="flex" direction="column">
                        <Typography component="h4" variant="body1" sx={{ fontSize: '1.0625rem' }} mb={1}>
                            Worker group: Python_1
                            <Typography ml={3} variant="subtitle1" style={{ display: 'inline' }}>
                                <FontAwesomeIcon icon={faDocker} style={{ marginRight: 4 }} />
                                Docker
                            </Typography>
                        </Typography>
                        <Typography component="h5" variant="subtitle1">
                            Python workers for generic work loads.
                        </Typography>
                    </Grid>

                    <Grid item display="flex" alignItems="center" sx={{ marginLeft: 'auto', marginRight: '2px' }}>
                        <div>
                            <Typography component="h2" variant="h2" align="right" sx={{ fontWeight: 900 }}>
                                6
                            </Typography>
                            <Typography variant="body1" sx={{ fontSize: '1.0625rem' }}>
                                Queue
                            </Typography>
                        </div>

                        <div style={{ marginLeft: 38 }}>
                            <Typography component="h2" variant="h2" align="right" sx={{ fontWeight: 900 }}>
                                6
                            </Typography>
                            <Typography variant="body1" sx={{ fontSize: '1.0625rem' }}>
                                Running
                            </Typography>
                        </div>

                        <div style={{ marginLeft: 38 }}>
                            <Typography component="h2" variant="h2" align="right" sx={{ fontWeight: 900 }}>
                                6
                            </Typography>
                            <Typography variant="body1" sx={{ fontSize: '1.0625rem' }}>
                                Succeeded
                            </Typography>
                        </div>

                        <div style={{ marginLeft: 38 }}>
                            <Typography component="h2" variant="h2" align="right" sx={{ fontWeight: 900 }}>
                                6
                            </Typography>
                            <Typography variant="body1" sx={{ fontSize: '1.0625rem' }}>
                                Failed
                            </Typography>
                        </div>
                    </Grid>
                </Grid>
            </Box>

            <Box component="table" mt={2} sx={{ width: tableWidth }} {...getTableProps()}>
                {/* <thead>
                    {headerGroups.map((headerGroup) => (
                        <Box
                            component="tr"
                            display="grid"
                            sx={{ '*:first-of-type': { ml: '22px' }, '*:last-child': { textAlign: 'center' } }}
                            gridTemplateColumns="repeat(3, 1fr)"
                            justifyContent="flex-start"
                            {...headerGroup.getHeaderGroupProps()}>
                            {headerGroup.headers.map((column) => (
                                <Box component="td" color="text.primary" fontWeight="600" fontSize="15px" textAlign="left" {...column.getHeaderProps()}>
                                    {column.render('Header')}
                                </Box>
                            ))}
                        </Box>
                    ))}
                </thead> */}
                <Box component="tbody" display="flex" sx={{ flexDirection: 'column' }} {...getTableBodyProps()}>
                    {rows.map((row, i) => {
                        prepareRow(row);
                        return (
                            <Box
                                component="tr"
                                {...row.getRowProps()}
                                display="grid"
                                gridTemplateColumns="repeat(3, 1fr)"
                                alignItems="start"
                                borderRadius="5px"
                                backgroundColor="background.secondary"
                                mt={2}
                                sx={{
                                    border: 1,
                                    borderColor: 'divider',
                                    padding: '15px 0',
                                    cursor: 'pointer',
                                    '&:hover': { background: 'background.hoverSecondary' },
                                    'td:last-child': { textAlign: 'center' },
                                }}>
                                {row.cells.map((cell) => {
                                    return (
                                        <Box component="td" {...cell.getCellProps()} textAlign="left">
                                            {cell.render('Cell')}
                                        </Box>
                                    );
                                })}
                            </Box>
                        );
                    })}
                </Box>
            </Box>
        </Box>
    );
}

const CustomWorker = ({ row }) => {
    const [id, status, lastUpdate, queue, running, succeeded, failed] = row.value;

    return (
        <Grid container direction="column" mx="22px" alignItems="left" justifyContent="flex-start">
            <div>
                <Typography component="h4" variant="h3" sx={{ display: 'inline' }}>
                    {id}
                </Typography>
                <Typography
                    component="h4"
                    variant="subtitle1"
                    color={status === 'Online' ? 'green' : 'red'}
                    fontWeight={700}
                    ml={3}
                    sx={{ display: 'inline', verticalAlign: 'top' }}>
                    {status}
                </Typography>
            </div>
            <Typography component="h5" mt={0.5} variant="subtitle1">
                {lastUpdate}
            </Typography>
            <Grid item display="flex" alignItems="center" mt={2} sx={{ alignSelf: 'flex-start' }}>
                <CustomChip amount={queue} label="Queue" margin={2} customColor="purple" />
                <CustomChip amount={running} label="Running" margin={1} customColor="orange" />
                <CustomChip amount={succeeded} label="Succeeded" margin={1} customColor="green" />
                <CustomChip amount={failed} label="Failed" margin={1} customColor="red" />
            </Grid>
        </Grid>
    );
};

// const CustomMemory = ({ row }) => {
//     const [percentage, mb] = row.value;

//     return (
//         <Grid container direction="column" alignItems="flex-start" flexDirection="row" pr={1}>
//             <Grid item>
//                 <Typography variant="h2" align="right" sx={{ fontWeight: 900 }}>
//                     {percentage}%
//                 </Typography>
//                 <Typography variant="h2" align="right" sx={{ fontWeight: 900 }}>
//                     {mb}MB
//                 </Typography>
//                 <Typography variant="body1" align="right" sx={{ fontSize: '1.0625rem' }}>
//                     Memory
//                 </Typography>
//             </Grid>
//             <img src={graph2} alt="" width="250px" style={{ marginLeft: '5px' }} />
//         </Grid>
//     );
// };

const dummyData = [
    {
        worker: {
            id: '36474-6768-6768-67859',
            name: 'Python 1',
            description: 'Python workers for generic work loads.',
            status: 'Online',
            type: 'Docker',
            lastUpdate: '2022-01-20T11:56:08Z',
            queue: 2,
            running: 2,
            succeeded: 2,
            failed: 2,
        },
        cpu: {
            percentage: 10,
            load: 0.27,
        },
        memory: {
            percentage: 10,
            mb: 200,
        },
    },
    {
        worker: {
            id: '36474-6768-6768-67859',
            name: 'Python 1',
            description: 'Python workers for generic work loads.',
            status: 'Online',
            type: 'Docker',
            lastUpdate: '2022-01-20T11:56:08Z',
            queue: 2,
            running: 2,
            succeeded: 2,
            failed: 2,
        },
        cpu: {
            percentage: 10,
            load: 0.27,
        },
        memory: {
            percentage: 10,
            mb: 200,
        },
    },
    {
        worker: {
            id: '36474-6768-6768-67859',
            name: 'Python 1',
            description: 'Python workers for generic work loads.',
            status: 'Online',
            type: 'Docker',
            lastUpdate: '2022-01-20T11:56:08Z',
            queue: 2,
            running: 2,
            succeeded: 2,
            failed: 2,
        },
        cpu: {
            percentage: 10,
            load: 0.27,
        },
        memory: {
            percentage: 10,
            mb: 200,
        },
    },
];
