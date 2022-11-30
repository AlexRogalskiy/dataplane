import { Box, Button, Drawer, Grid, Tooltip, Typography } from '@mui/material';
import React, { useEffect, useMemo, useState } from 'react';
import { useGlobalFilter, useTable } from 'react-table';
import CustomChip from '../../../components/CustomChip';
import Search from '../../../components/Search';
import { faEdit, faFilter } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import AddRPAWorkerDrawer from '../../../components/DrawerContent/AddRPAWorker';
import EditRPAWorkerDrawer from '../../../components/DrawerContent/EditRPAWorker';
import { useHistory } from 'react-router-dom';
import { useGlobalEnvironmentState } from '../../../components/EnviromentDropdown';
import { useGetRemoteWorkers } from '../../../graphql/getRemoteWorkers';
import { useSnackbar } from 'notistack';
import ConnectRemoteWorkerDrawer from '../../../components/DrawerContent/ConnectRemoteWorkerDrawer';

const tableWidth = '850px';

export default function RPAWorkers() {
    const [showAddWorkerDrawer, setShowAddWorkerDrawer] = useState(false);
    const [showEditWorkerDrawer, setShowEditWorkerDrawer] = useState(false);
    const [showConnectDrawer, setShowConnectDrawer] = useState(false);
    const [selectedWorker, setSelectedWorker] = useState(null);
    const [remoteWorkers, setRemoteWorkers] = useState([]);

    const history = useHistory();

    // Global environment state with hookstate
    const Environment = useGlobalEnvironmentState();

    // Graphql hook
    const getRemoteWorkers = useGetRemoteWorkersHook(Environment.id.get(), setRemoteWorkers);

    useEffect(() => {
        getRemoteWorkers();

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [Environment.id.get()]);

    const columns = useMemo(
        () => [
            {
                Header: 'Worker',
                accessor: (row) => [row.WorkerName, row.Status, row.WorkerID],
                Cell: (row) => (
                    <Box display="flex" flexDirection="column">
                        <Box display="flex">
                            <Tooltip title={row.value[2]} placement="top">
                                <Typography variant="caption" lineHeight={1.2} mr={0.5}>
                                    {row.value[0]}
                                </Typography>
                            </Tooltip>
                            <Box position="relative">
                                <Box
                                    onClick={() => {
                                        setShowEditWorkerDrawer(true);
                                        setSelectedWorker(row.value);
                                    }}
                                    component={FontAwesomeIcon}
                                    fontSize={12}
                                    color="#7D7D7D"
                                    icon={faEdit}
                                    cursor="pointer"
                                    position="absolute"
                                    left="5px"
                                    top="1px"
                                />
                            </Box>
                        </Box>

                        <Typography variant="caption" lineHeight={1.2} fontWeight={700} color={row.value[1] === 'online' ? 'success.main' : 'red'}>
                            {row.value[1]}
                        </Typography>
                    </Box>
                ),
            },
            {
                Header: 'Process groups',
                accessor: 'groupCount',
                Cell: (row) => (
                    <Typography mt={-2} variant="caption" color="cyan.main" sx={{ cursor: 'pointer' }}>
                        Manage(1)
                    </Typography>
                ),
            },
            {
                Header: 'Last ping',
                accessor: 'lastPing',
                Cell: (row) => (
                    <Typography mt={-2} variant="caption">
                        {row.LastPing}
                    </Typography>
                ),
            },
            {
                Header: 'Manage',
                accessor: (row) => [row.WorkerID, row.WorkerName],
                Cell: (row) => {
                    return (
                        <>
                            <Typography
                                variant="caption"
                                mr={1}
                                mt={-2}
                                color="cyan.main"
                                sx={{ cursor: 'pointer' }}
                                onClick={() => history.push(`/remote/workers/${row.value[0]}`)}>
                                Configure
                            </Typography>
                            <Typography variant="caption" mt={-2}>
                                |
                            </Typography>
                            <Typography
                                variant="caption"
                                ml={1}
                                mt={-2}
                                color="cyan.main"
                                sx={{ cursor: 'pointer' }}
                                onClick={() => {
                                    setSelectedWorker(row.value);
                                    setShowConnectDrawer(true);
                                }}>
                                Connect
                            </Typography>
                        </>
                    );
                },
            },
        ],

        []
    );

    // Use the state and functions returned from useTable to build your UI
    const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow, setGlobalFilter } = useTable(
        {
            columns,
            data: remoteWorkers,
        },

        useGlobalFilter
    );

    return (
        <Box className="page">
            <Box display="flex" alignItems="center" width={tableWidth}>
                <Typography id="test" component="h2" variant="h2" color="text.primary">
                    RPA Workers
                </Typography>

                <Button onClick={() => history.push('/remote/processgroups')} variant="text" sx={{ marginLeft: 'auto', marginRight: 2 }}>
                    Manage process groups
                </Button>

                <Button variant="contained" size="small" onClick={() => setShowAddWorkerDrawer(true)}>
                    Add worker
                </Button>
            </Box>

            <Box mt={'45px'} sx={{ width: tableWidth }}>
                <Grid container mt={4} direction="row" alignItems="center" justifyContent="flex-start">
                    <Grid item display="flex" alignItems="center" sx={{ alignSelf: 'center' }}>
                        <CustomChip amount={remoteWorkers.length} label="RPA Workers" margin={2} customColor="orange" />
                    </Grid>

                    <Grid item display="flex" alignItems="center" sx={{ alignSelf: 'center' }}>
                        <Box component={FontAwesomeIcon} icon={faFilter} sx={{ fontSize: 12, color: '#b9b9b9' }} />

                        <Typography variant="subtitle1" color="#737373" ml={1}>
                            Process group = Python 1
                        </Typography>
                    </Grid>

                    <Grid item display="flex" alignItems="center" sx={{ marginLeft: 'auto', marginRight: '2px' }}>
                        <Search placeholder="Find workers" onChange={setGlobalFilter} width="290px" />
                    </Grid>
                </Grid>
            </Box>

            {remoteWorkers.length > 0 ? (
                <Box component="table" mt={6} sx={{ width: tableWidth }} {...getTableProps()}>
                    <Box component="thead" display="flex" sx={{ flexDirection: 'column' }}>
                        {headerGroups.map((headerGroup) => (
                            <Box //
                                component="tr"
                                {...headerGroup.getHeaderGroupProps()}
                                display="grid"
                                gridTemplateColumns="2fr repeat(3, 1fr)"
                                textAlign="left">
                                {headerGroup.headers.map((column, idx) => (
                                    <Box
                                        component="th"
                                        {...column.getHeaderProps()}
                                        sx={{
                                            borderTopLeftRadius: idx === 0 ? 5 : 0,
                                            borderTopRightRadius: idx === headerGroup.headers.length - 1 ? 5 : 0,
                                            border: '1px solid',
                                            borderLeft: idx === 0 ? '1px solid' : 0,
                                            borderColor: 'editorPage.borderColor',
                                            backgroundColor: 'background.worker',
                                            pl: 1,
                                            fontSize: '0.875rem',
                                            py: '6px',
                                        }}>
                                        {column.render('Header')}
                                    </Box>
                                ))}
                            </Box>
                        ))}
                    </Box>
                    <Box component="tbody" display="flex" sx={{ flexDirection: 'column' }} {...getTableBodyProps()}>
                        {rows.map((row, i) => {
                            prepareRow(row);
                            return (
                                <Box
                                    component="tr"
                                    {...row.getRowProps()}
                                    display="grid"
                                    gridTemplateColumns="2fr repeat(3, 1fr)"
                                    alignItems="start"
                                    backgroundColor="background.secondary">
                                    {row.cells.map((cell, idx) => {
                                        return (
                                            <Box
                                                component="td"
                                                {...cell.getCellProps()}
                                                sx={{
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    border: '1px solid',
                                                    borderColor: 'editorPage.borderColor',
                                                    height: '50px',
                                                    // If first cell
                                                    borderLeft: idx === 0 ? '1px solid editorPage.borderColor' : 0,
                                                    // If last cell of last row
                                                    borderBottomRightRadius: idx === row.cells.length - 1 && i === rows.length - 1 ? '5px' : 0,
                                                    // If first cell of last row
                                                    borderBottomLeftRadius: idx === 0 && i === rows.length - 1 ? '5px' : 0,
                                                    borderTop: 0,
                                                    pl: 1,
                                                    textAlign: 'left',
                                                    // If last cell
                                                    ...(idx === row.cells.length - 1 && {
                                                        justifyContent: 'center',
                                                        paddingLeft: '0',
                                                    }),
                                                }}>
                                                {cell.render('Cell')}
                                            </Box>
                                        );
                                    })}
                                </Box>
                            );
                        })}
                    </Box>
                </Box>
            ) : null}

            {/* Add worker drawer */}
            <Drawer anchor="right" open={showAddWorkerDrawer} onClose={() => setShowAddWorkerDrawer(!showAddWorkerDrawer)}>
                <AddRPAWorkerDrawer
                    handleClose={() => {
                        setShowAddWorkerDrawer(false);
                    }}
                    getRemoteWorkers={getRemoteWorkers}
                />
            </Drawer>

            {/* Edit worker drawer */}
            <Drawer anchor="right" open={showEditWorkerDrawer} onClose={() => setShowEditWorkerDrawer(!showEditWorkerDrawer)}>
                <EditRPAWorkerDrawer
                    handleClose={() => {
                        setShowEditWorkerDrawer(false);
                    }}
                />
            </Drawer>

            {/* Connect drawer */}
            <Drawer
                //
                hideBackdrop
                sx={{ width: 'calc(100% - 203px)', zIndex: 1099, [`& .MuiDrawer-paper`]: { width: 'calc(100% - 203px)', top: 82 } }}
                anchor="right"
                open={showConnectDrawer}
                onClose={() => setShowConnectDrawer(!false)}>
                <ConnectRemoteWorkerDrawer
                    handleClose={() => {
                        setShowConnectDrawer(false);
                    }}
                    worker={selectedWorker}
                    environmentID={Environment.id.get()}
                />
            </Drawer>
        </Box>
    );
}

// ** Custom Hooks
const useGetRemoteWorkersHook = (environmentID, setRemoteWorkers) => {
    // GraphQL hook
    const getRemoteWorkers = useGetRemoteWorkers();

    const { enqueueSnackbar } = useSnackbar();

    // Get worker groups
    return async () => {
        const response = await getRemoteWorkers({ environmentID });

        if (response === null) {
            setRemoteWorkers([]);
        } else if (response.r || response.error) {
            enqueueSnackbar("Can't get remote process groups: " + (response.msg || response.r || response.error), { variant: 'error' });
        } else if (response.errors) {
            response.errors.map((err) => enqueueSnackbar(err.message, { variant: 'error' }));
        } else {
            setRemoteWorkers(response);
        }
    };
};
