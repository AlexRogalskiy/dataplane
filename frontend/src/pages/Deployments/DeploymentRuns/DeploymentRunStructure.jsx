import { useSnackbar } from "notistack";
import { useGetSinglepipelineRunAndTasks } from "../../../graphql/getSinglepipelineRunAndTasks";
import { useGlobalPipelineRun } from "../../PipelineRuns/GlobalPipelineRunUIState";


/* 
Get the structure for a pipeline run
*/
export const GetDeploymentRun = () => {

    const { enqueueSnackbar } = useSnackbar();
    const FlowState = useGlobalPipelineRun();
    const getPipelineRun = useGetSinglepipelineRunAndTasks();

    return async (pipelineID, runID, environmentID) => {

                // Get single pipelines run and statuses
                let [singleRunResponse, tasksResponse] = await getPipelineRun({
                    pipelineID: pipelineID,
                    runID: runID,
                    environmentID: environmentID,
                });

                if (singleRunResponse.length === 0) {
                    // setSelectedRun([]);
                } else if (singleRunResponse.r || singleRunResponse.error) {
                    enqueueSnackbar("Can't get pipeline run: " + (singleRunResponse.msg || singleRunResponse.r || singleRunResponse.error), { variant: 'error' });
                } else if (singleRunResponse.errors) {
                    singleRunResponse.errors.map((err) => enqueueSnackbar(err.message, { variant: 'error' }));
                } else {
                    // setSelectedRun(singleRunResponse);
                    FlowState.elements.set(singleRunResponse.run_json);
 
                }
            }
            };