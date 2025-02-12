import React, {useState} from 'react';
import {Dimension, migrateLegacyQuery, VariableQuery, VariableQueryType} from './types';
import {DataSource} from './datasource';
import {SelectableValue} from '@grafana/data'
import DimensionSettings from './components/DimensionSettings';
import {AsyncSelect, InlineField, Input, Select} from '@grafana/ui';

const formatDefinition = (query: VariableQuery): string => {
    switch (query.queryType) {
        case VariableQueryType.metric:
            return query.dimensions.map((x) => `${x.key}=${x.value}`).join(';');
        case VariableQueryType.dimensionValue:
            return `dimension=${query.dimensionKey}&filter=${query.dimensionValueFilter}`;
    }
};
const VariableQueryEditor = (props: {
    query: VariableQuery | string;
    onChange: (query: VariableQuery, definition: string) => void;
    datasource: DataSource;
}) => {
    const {datasource, onChange} = props;
    const query = migrateLegacyQuery(props.query);
    const [state, updateState] = useState(query);

    const onChangeQueryType = (qt?: VariableQueryType) => {
        const newState = {...state, queryType: qt || VariableQueryType.metric};
        updateState(newState);
        onChange(newState, formatDefinition(newState));
    };

    const onDimensionsChange = (dimensions: Dimension[]) => {
        const newState = {...state, dimensions: dimensions};
        updateState(newState);
        onChange(newState, formatDefinition(newState));
    };

    const onDimensionKeyChange = (key?: string) => {
        const newState = {...state, dimensionKey: key || ''};
        updateState(newState);
        onChange(newState, formatDefinition(newState));
    };
    const loadDimensionKeys = (query: string): Promise<Array<SelectableValue<string>>> => {
        return datasource.listDimensionKeys(query, []);
    };
    const options: Array<SelectableValue<VariableQueryType>> = [
        {
            value: VariableQueryType.metric,
            label: 'Metric',
            description: 'the query selects metrics',
        },
        {
            value: VariableQueryType.dimensionValue,
            label: 'DimensionValue',
            description: 'the query selects dimension values',
        },
    ];

    function onDimensionValueFilterChange(filter: string) {
        const newState = {...state, dimensionValueFilter: filter}
        updateState(newState);
        onChange(newState, formatDefinition(newState));
    }

    return (
        <>
            <div className="gf-form">
                <label className="gf-form-label width-10">Query Type</label>
                <Select onChange={(x) => onChangeQueryType(x.value)} options={options} value={state.queryType}/>
            </div>
            {state.queryType === VariableQueryType.metric && (
                <DimensionSettings initState={state.dimensions} onChange={onDimensionsChange} datasource={datasource}/>
            )}
            {state.queryType === VariableQueryType.dimensionValue && (
                <>
                    <div className="gf-form">
                        <label className="gf-form-label width-10">Dimension Key</label>
                        <AsyncSelect
                            defaultOptions={true}
                            value={{label: state.dimensionKey, value: state.dimensionKey}}
                            cacheOptions={false}
                            loadOptions={loadDimensionKeys}
                            onChange={(e) => onDimensionKeyChange(e.value)}
                        />
                        <InlineField label="Filter" labelWidth={20}
                                     tooltip={"filter dimension values"}>
                            <Input width={40} onChange={x => onDimensionValueFilterChange(x.currentTarget.value)}
                                   value={state.dimensionValueFilter}/>
                        </InlineField>
                    </div>
                </>
            )}
        </>
    );
};

export default VariableQueryEditor;
